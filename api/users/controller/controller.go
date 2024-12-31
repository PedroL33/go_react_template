package controller

import (
	"context"
	"example/dashboard/api/db"
	"example/dashboard/api/models"
	"example/dashboard/api/users"
	"example/dashboard/api/users/payloads"
	"example/dashboard/config"
	http_errors "example/dashboard/errors"
	"example/dashboard/util"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
)

type usersController struct {
	cfg        *config.Config
	usersStore users.Store
	txnManager db.TransactionManager
	logger     util.Logger
}

func NewUsersController(cfg *config.Config, usersStore users.Store, txnManager db.TransactionManager, logger util.Logger) users.Controller {
	return &usersController{cfg: cfg, usersStore: usersStore, txnManager: txnManager, logger: logger}
}

func (uc *usersController) CreateUser(ctx context.Context, request *payloads.CreateUserRequest) (*models.UserWithToken, error) {
	var err error
	var existingUser *models.User
	if existingUser, err = uc.usersStore.GetUserByEmail(ctx, request.Email, nil); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, http_errors.NewInternalServerError(
				errors.Wrap(err, "UsersController.CreateUser"),
				"Internal sever error.",
			)
		}
	}

	if existingUser != nil {
		return nil, http_errors.NewHttpError(
			http.StatusBadRequest,
			"Email already exists.",
			errors.Wrap(err, "UsersController.CreateUser"),
		)
	}

	user := &models.User{
		Email:    request.Email,
		Password: request.Password,
	}
	user.PrepareCreate()

	var tx db.DbConn
	if tx, err = uc.txnManager.Begin(ctx); err != nil {
		return nil, http_errors.NewInternalServerError(
			errors.Wrap(err, "UsersController.CreateUser"),
			"Internal server error. Failed to begin transaction.",
		)
	}
	defer tx.Rollback(ctx)

	var createdUser *models.User
	if createdUser, err = uc.usersStore.CreateUser(ctx, user, tx); err != nil {
		return nil, http_errors.NewBadRequestError(
			errors.Wrap(err, "UsersController.CreateUser"),
			"User creation failed.",
		)
	}

	var token string
	if token, err = util.CreateToken(uc.cfg, createdUser); err != nil {
		return nil, http_errors.NewInternalServerError(
			errors.Wrap(err, "UsersController.CreateUser"),
			"Internal server error. Failed to generate token.",
		)
	}

	createdUser.Sanitize()

	if err = tx.Commit(ctx); err != nil {
		return nil, http_errors.NewInternalServerError(
			errors.Wrap(err, "UsersController.CreateUser"),
			"Internal server error. Failed to commit transaction.",
		)
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

func (uc *usersController) Login(ctx context.Context, user *models.User) (*payloads.LoginResponse, error) {
	var foundUser *models.User
	var err error
	if foundUser, err = uc.usersStore.GetUserByEmail(ctx, user.Email, nil); err != nil {
		return nil, http_errors.NewForbiddenError(
			errors.Wrap(err, "UsersController.Login"),
			"Invalid credentials.",
		)
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, http_errors.NewForbiddenError(
			errors.Wrap(err, "UsersController.Login"),
			"Invalid credentials.",
		)
	}

	if foundUser.IsTwoFactorEnabled.Bool {
		if err = uc.usersStore.DeleteLoginSessionByUserId(ctx, foundUser.Id, nil); err != nil {
			return nil, http_errors.NewInternalServerError(
				errors.Wrap(err, "UsersController.Login"),
				"Two factor authentication error.",
			)
		}
		var loginSessionId int
		if loginSessionId, err = uc.usersStore.CreateLoginSession(ctx, foundUser.Id, nil); err != nil {
			return nil, http_errors.NewInternalServerError(
				errors.Wrap(err, "UsersController.Login"),
				"Two factor authentication error.",
			)
		}

		return &payloads.LoginResponse{
			LoginSessionId: loginSessionId,
		}, nil
	} else {
		foundUser.Sanitize()
		var token string
		if token, err = util.CreateToken(uc.cfg, foundUser); err != nil {
			return nil, http_errors.NewInternalServerError(
				errors.Wrap(err, "UsersController.Login"),
				"Error while creating token.",
			)
		}

		return &payloads.LoginResponse{
			User:  foundUser,
			Token: token,
		}, nil
	}
}

func (uc *usersController) VerifyLogin(ctx context.Context, verifyLoginRequest *payloads.LoginWithOptCodeRequest) (*payloads.LoginResponse, error) {
	var err error

	var loginSession = &models.LoginSession{}
	if loginSession, err = uc.usersStore.GetLoginSessionById(ctx, verifyLoginRequest.LoginSessionId, nil); err != nil {
		return nil, http_errors.NewBadRequestError(errors.Wrap(err, "UsersController.VerifyLogin"), "Unable to locate login session.")
	}

	if loginSession.Expiration.Before(time.Now()) {
		return nil, http_errors.NewExpiredSessionError(errors.Wrap(err, "UsersController.VerifyLogin"), "Expired session.")
	}

	var currentUser = &models.User{}
	if currentUser, err = uc.usersStore.GetUserById(ctx, loginSession.UserId, nil); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.VerifyLogin"), "Something went wrong.")
	}

	var isOptCodeValid bool
	if isOptCodeValid = totp.Validate(verifyLoginRequest.OtpCode, currentUser.TwoFactorSecret.String); !isOptCodeValid {
		return nil, http_errors.NewBadRequestError(
			errors.Wrap(err, "UsersController.VerifyLogin"),
			"Code not valid.",
		)
	}

	currentUser.Sanitize()
	var token string
	if token, err = util.CreateToken(uc.cfg, currentUser); err != nil {
		return nil, http_errors.NewInternalServerError(
			errors.Wrap(err, "UsersController.VerifyLogin"),
			"Error while creating token.",
		)
	}

	return &payloads.LoginResponse{
		User:  currentUser,
		Token: token,
	}, nil
}

func (uc *usersController) VerifyLoginWithRecoveryCode(ctx context.Context, verifyLoginRequest *payloads.LoginWithRecoveryCodeRequest) (*payloads.LoginResponse, error) {
	var err error

	var loginSession = &models.LoginSession{}
	if loginSession, err = uc.usersStore.GetLoginSessionById(ctx, verifyLoginRequest.LoginSessionId, nil); err != nil {
		return nil, http_errors.NewBadRequestError(errors.Wrap(err, "UsersController.VerifyLoginWithRecoveryCode"), "Unable to locate login session.")
	}

	if loginSession.Expiration.Before(time.Now()) {
		return nil, http_errors.NewExpiredSessionError(errors.Wrap(err, "UsersController.VerifyLoginWithRecoveryCode"), "Expired session.")
	}

	var currentUser = &models.User{}
	if currentUser, err = uc.usersStore.GetUserById(ctx, loginSession.UserId, nil); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.VerifyLoginWithRecoveryCode"), "Something went wrong.")
	}

	var recoveryCodes []*models.RecoveryCode
	if recoveryCodes, err = uc.usersStore.GetRecoveryCodesByUserId(ctx, currentUser.Id, nil); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.VerifyLoginWithRecoveryCode"), "Error while verifying recovery codes.")
	}

	for i := 0; i < len(recoveryCodes); i++ {
		if recoveryCodes[i].Code == verifyLoginRequest.RecoveryCode && !recoveryCodes[i].IsRedeemed.Bool {
			currentUser.Sanitize()

			if err = uc.usersStore.RedeemRecoveryCode(ctx, recoveryCodes[i].Id, nil); err != nil {
				return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.VerifyLoginWithRecoveryCode"), "Error while redeeming recovery codes.")
			}

			var token string
			if token, err = util.CreateToken(uc.cfg, currentUser); err != nil {
				return nil, http_errors.NewInternalServerError(
					errors.Wrap(err, "UsersController.VerifyLoginWithRecoveryCode"),
					"Error while creating token.",
				)
			}

			return &payloads.LoginResponse{
				User:  currentUser,
				Token: token,
			}, nil
		}
	}

	return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.VerifyLoginWithRecoveryCode"), "Invalid recovery code.")
}

func (uc *usersController) Begin2faSetupSession(ctx context.Context, currentUser *models.User) (string, error) {
	var err error

	if currentUser.IsTwoFactorEnabled.Bool {
		return "", http_errors.NewBadRequestError(errors.Wrap(err, "UsersController.Begin2faSetupSession"), "Already enrolled in two factor authentication.")
	}

	twoFactorSession := &models.TwoFactorSetupSession{
		UserId:     currentUser.Id,
		Expiration: time.Now().Add(2 * time.Hour),
	}

	var qrCodeString string
	if qrCodeString, err = twoFactorSession.PopulateSecretStringAndReturnBase64QrCode(currentUser.Email); err != nil {
		return "", http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Begin2faSetupSession"), "Internal server error. Could not populate secret string.")
	}

	var tx db.DbConn
	if tx, err = uc.txnManager.Begin(ctx); err != nil {
		return "", http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Begin2faSetupSession"), "Internal server error. Could not begin transaction.")
	}
	defer tx.Rollback(ctx)

	if err = uc.usersStore.Delete2faSetupSession(ctx, currentUser.Id, tx); err != nil {
		return "", http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Begin2faSetupSession"), "Error deleting existing session.")
	}

	if _, err = uc.usersStore.Create2faSetupSession(ctx, twoFactorSession, tx); err != nil {
		return "", http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Begin2faSetupSession"), "Internal server error. Could not create 2fa session.")
	}

	if err = tx.Commit(ctx); err != nil {
		return "", http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Begin2faSetupSession"), "Error while commiting transaction.")
	}

	return qrCodeString, nil
}

func (uc *usersController) Complete2faSetup(ctx context.Context, complete2faSetupRequest *payloads.Complete2faSetupRequest, currentUser *models.User) ([]*models.RecoveryCode, error) {
	var err error

	var tx db.DbConn
	if tx, err = uc.txnManager.Begin(ctx); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Complete2faSetup"), "Internal server error. Could not begin transaction.")
	}
	defer tx.Rollback(ctx)

	var currentUserSetupSession *models.TwoFactorSetupSession
	if currentUserSetupSession, err = uc.usersStore.Get2faSetupSessionByUserId(ctx, currentUser.Id, nil); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Complete2faSetup"), "Could not find related two facto setup session.")
	}

	var isOptCodeValid bool
	if isOptCodeValid = totp.Validate(complete2faSetupRequest.OtpCode, currentUserSetupSession.SecretString); !isOptCodeValid {
		return nil, http_errors.NewBadRequestError(
			errors.Wrap(err, "UsersController.Complete2faSetup"),
			"Code not valid.",
		)
	}

	if err = uc.usersStore.EnableTwoFactorAuth(ctx, currentUserSetupSession, tx); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Complete2faSetup"), "Error while enabling two factor authentication.")
	}

	recoveryCodes := make([]*models.RecoveryCode, 0, 10)
	for i := 0; i < 10; i++ {
		code, err := uc.usersStore.GenerateRecoveryCode(ctx, currentUserSetupSession.UserId, tx)
		if err != nil {
			return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Complete2faSetup"), "Error while generating recovery codes.")
		}
		recoveryCodes = append(recoveryCodes, code)
	}

	if err = uc.usersStore.Delete2faSetupSession(ctx, currentUser.Id, tx); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Complete2faSetup"), "Error deleting setup session.")
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Complete2faSetup"), "Error while commiting transaction.")
	}

	return recoveryCodes, nil
}

func (uc *usersController) Disable2fa(ctx context.Context, currentUser *models.User, disable2faRequest *payloads.Disable2faRequest) error {
	var tx db.DbConn
	var err error
	if tx, err = uc.txnManager.Begin(ctx); err != nil {
		return http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Disable2fa"), "Internal server error. Could not begin transaction.")
	}
	defer tx.Rollback(ctx)

	if err = currentUser.ComparePasswords(disable2faRequest.Password); err != nil {
		return http_errors.NewForbiddenError(
			errors.Wrap(err, "UsersController.Disable2fa"),
			"Invalid credentials.",
		)
	}

	if err = uc.usersStore.DisableTwoFactorAuth(ctx, currentUser.Id, tx); err != nil {
		return http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Disable2fa"), "Error disabling two factor authentication.")
	}

	if err = uc.usersStore.DeleteRecoveryCodes(ctx, currentUser.Id, tx); err != nil {
		return http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Disable2fa"), "Error deleting recovery codes.")
	}

	if err = tx.Commit(ctx); err != nil {
		return http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Disable2fa"), "Error commiting transaction.")
	}

	return nil
}

func (uc *usersController) UpdatePassword(ctx context.Context, currentUser *models.User, updatePasswordRequest *payloads.UpdatePasswordRequest) error {
	var err error

	if err = currentUser.ComparePasswords(updatePasswordRequest.CurrentPassword); err != nil {
		return http_errors.NewForbiddenError(
			errors.Wrap(err, "UsersController.UpdatePassword"),
			"Invalid credentials.",
		)
	}

	currentUser.Password = updatePasswordRequest.NewPassword
	if err := currentUser.HashPassword(); err != nil {
		return http_errors.NewInternalServerError(
			errors.Wrap(err, "UsersController.UpdatePassword"),
			"Error while hashing password.",
		)
	}

	if err = uc.usersStore.UpdatePassword(ctx, currentUser.Password, currentUser.Id, nil); err != nil {
		return http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.UpdatePassword"), "Error while updating password.")
	}

	return nil
}

func (uc *usersController) RegenerateRecoveryCodes(ctx context.Context, currentUser *models.User) ([]*models.RecoveryCode, error) {
	var err error
	var tx db.DbConn
	if tx, err = uc.txnManager.Begin(ctx); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.RegenerateRecoveryCodes"), "Internal server error. Could not begin transaction.")
	}
	defer tx.Rollback(ctx)

	if err = uc.usersStore.DeleteRecoveryCodes(ctx, currentUser.Id, tx); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.RegenerateRecoveryCodes"), "Error deleting recovery codes.")
	}

	recoveryCodes := make([]*models.RecoveryCode, 0, 10)
	for i := 0; i < 10; i++ {
		code, err := uc.usersStore.GenerateRecoveryCode(ctx, currentUser.Id, tx)
		if err != nil {
			return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.RegenerateRecoveryCodes"), "Error generating recovery code.")
		}
		recoveryCodes = append(recoveryCodes, code)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.RegenerateRecoveryCodes"), "Error while commiting transaction.")
	}

	return recoveryCodes, nil
}
