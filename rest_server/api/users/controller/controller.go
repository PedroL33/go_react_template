package controller

import (
	"context"
	"example/template/rest_server/api/db"
	"example/template/rest_server/api/httpcomm"
	"example/template/rest_server/api/models"
	"example/template/rest_server/api/users"
	"example/template/rest_server/api/users/payloads"
	"example/template/rest_server/config"
	"example/template/rest_server/logger"
	"example/template/rest_server/util"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
)

type usersController struct {
	cfg        *config.AppConfig
	usersStore users.Store
	txnManager db.TransactionManager
	logger     logger.Logger
}

func NewUsersController(cfg *config.AppConfig, usersStore users.Store, txnManager db.TransactionManager, logger logger.Logger) users.Controller {
	return &usersController{cfg: cfg, usersStore: usersStore, txnManager: txnManager, logger: logger}
}

func (uc *usersController) CreateUser(ctx context.Context, request *payloads.CreateUserRequest) (*models.UserWithToken, error) {
	existingUser, err := uc.usersStore.GetUserByUsername(ctx, request.Username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Internal sever error.")
	}

	if existingUser != nil {
		return nil, httpcomm.NewBadRequestError(util.Wrap(err), "Username already exists.")
	}

	user := &models.User{
		Username: request.Username,
		Password: request.Password,
	}
	user.PrepareCreate()

	var createdUser *models.User
	err = uc.txnManager.WithTx(ctx, func(tx db.Querier) error {
		var txErr error
		createdUser, txErr = uc.usersStore.WithQuerier(tx).CreateUser(ctx, user)
		return txErr
	})
	if err != nil {
		return nil, httpcomm.NewBadRequestError(util.Wrap(err), "User creation failed.")
	}

	token, err := util.CreateToken(uc.cfg, createdUser)
	if err != nil {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Internal server error. Failed to generate token.")
	}

	createdUser.Sanitize()

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

func (uc *usersController) Login(ctx context.Context, user *models.User) (*payloads.LoginResponse, error) {
	foundUser, err := uc.usersStore.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, httpcomm.NewForbiddenError(util.Wrap(err), "Invalid credentials.")
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, httpcomm.NewForbiddenError(util.Wrap(err), "Invalid credentials.")
	}

	if foundUser.IsTwoFactorEnabled.Bool {
		if err = uc.usersStore.DeleteLoginSessionByUserId(ctx, foundUser.Id); err != nil {
			return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Two factor authentication error.")
		}
		loginSessionId, err := uc.usersStore.CreateLoginSession(ctx, foundUser.Id)
		if err != nil {
			return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Two factor authentication error.")
		}

		return &payloads.LoginResponse{
			LoginSessionId: loginSessionId,
		}, nil
	}

	foundUser.Sanitize()
	token, err := util.CreateToken(uc.cfg, foundUser)
	if err != nil {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Error while creating token.")
	}

	return &payloads.LoginResponse{
		User:  foundUser,
		Token: token,
	}, nil
}

func (uc *usersController) VerifyLogin(ctx context.Context, verifyLoginRequest *payloads.LoginWithOptCodeRequest) (*payloads.LoginResponse, error) {
	loginSession, err := uc.usersStore.GetLoginSessionById(ctx, verifyLoginRequest.LoginSessionId)
	if err != nil {
		return nil, httpcomm.NewBadRequestError(util.Wrap(err), "Unable to locate login session.")
	}

	if loginSession.Expiration.Before(time.Now()) {
		return nil, httpcomm.NewUnauthorizedError(util.Wrap(err), "Expired session.")
	}

	currentUser, err := uc.usersStore.GetUserById(ctx, loginSession.UserId)
	if err != nil {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Something went wrong.")
	}

	if !totp.Validate(verifyLoginRequest.OtpCode, currentUser.TwoFactorSecret.String) {
		return nil, httpcomm.NewBadRequestError(util.Wrap(err), "Code not valid.")
	}

	currentUser.Sanitize()
	token, err := util.CreateToken(uc.cfg, currentUser)
	if err != nil {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Error while creating token.")
	}

	return &payloads.LoginResponse{
		User:  currentUser,
		Token: token,
	}, nil
}

func (uc *usersController) VerifyLoginWithRecoveryCode(ctx context.Context, verifyLoginRequest *payloads.LoginWithRecoveryCodeRequest) (*payloads.LoginResponse, error) {
	loginSession, err := uc.usersStore.GetLoginSessionById(ctx, verifyLoginRequest.LoginSessionId)
	if err != nil {
		return nil, httpcomm.NewBadRequestError(util.Wrap(err), "Unable to locate login session.")
	}

	if loginSession.Expiration.Before(time.Now()) {
		return nil, httpcomm.NewUnauthorizedError(util.Wrap(err), "Expired session.")
	}

	currentUser, err := uc.usersStore.GetUserById(ctx, loginSession.UserId)
	if err != nil {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Something went wrong.")
	}

	recoveryCodes, err := uc.usersStore.GetRecoveryCodesByUserId(ctx, currentUser.Id)
	if err != nil {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Error while verifying recovery codes.")
	}

	for i := 0; i < len(recoveryCodes); i++ {
		if recoveryCodes[i].Code == verifyLoginRequest.RecoveryCode && !recoveryCodes[i].IsRedeemed.Bool {
			currentUser.Sanitize()

			if err = uc.usersStore.RedeemRecoveryCode(ctx, recoveryCodes[i].Id); err != nil {
				return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Error while redeeming recovery codes.")
			}

			token, err := util.CreateToken(uc.cfg, currentUser)
			if err != nil {
				return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Error while creating token.")
			}

			return &payloads.LoginResponse{
				User:  currentUser,
				Token: token,
			}, nil
		}
	}

	return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Invalid recovery code.")
}

func (uc *usersController) Begin2faSetupSession(ctx context.Context, currentUser *models.User) (string, error) {
	if currentUser.IsTwoFactorEnabled.Bool {
		return "", httpcomm.NewBadRequestError(nil, "Already enrolled in two factor authentication.")
	}

	twoFactorSession := &models.TwoFactorSetupSession{
		UserId:     currentUser.Id,
		Expiration: time.Now().Add(2 * time.Hour),
	}

	qrCodeString, err := twoFactorSession.PopulateSecretStringAndReturnBase64QrCode(currentUser.Username)
	if err != nil {
		return "", httpcomm.NewInternalServerError(util.Wrap(err), "Internal server error. Could not populate secret string.")
	}

	err = uc.txnManager.WithTx(ctx, func(tx db.Querier) error {
		store := uc.usersStore.WithQuerier(tx)

		if err := store.Delete2faSetupSession(ctx, currentUser.Id); err != nil {
			// return httpcomm.NewInternalServerError(util.Wrap(err), "delete existing session")
			return util.Wrap(err)
		}
		if _, err := store.Create2faSetupSession(ctx, twoFactorSession); err != nil {
			// return httpcomm.NewInternalServerError(util.Wrap(err), "create 2fa session")
			return util.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return "", httpcomm.NewInternalServerError(util.Wrap(err), "Error while setting up two factor authentication.")
	}

	return qrCodeString, nil
}

func (uc *usersController) Complete2faSetup(ctx context.Context, complete2faSetupRequest *payloads.Complete2faSetupRequest, currentUser *models.User) ([]*models.RecoveryCode, error) {
	currentUserSetupSession, err := uc.usersStore.Get2faSetupSessionByUserId(ctx, currentUser.Id)
	if err != nil {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Could not find related two facto setup session.")
	}

	if !totp.Validate(complete2faSetupRequest.OtpCode, currentUserSetupSession.SecretString) {
		return nil, httpcomm.NewBadRequestError(nil, "Code not valid.")
	}

	recoveryCodes := make([]*models.RecoveryCode, 0, 10)
	err = uc.txnManager.WithTx(ctx, func(tx db.Querier) error {
		store := uc.usersStore.WithQuerier(tx)

		if err := store.EnableTwoFactorAuth(ctx, currentUserSetupSession); err != nil {
			return httpcomm.NewInternalServerError(util.Wrap(err), "Error while enabling two factor auth")
		}

		recoveryCodes = recoveryCodes[:0]
		for i := 0; i < 10; i++ {
			code, err := store.GenerateRecoveryCode(ctx, currentUserSetupSession.UserId)
			if err != nil {
				return httpcomm.NewInternalServerError(util.Wrap(err), "Error while generating recovery code")
			}
			recoveryCodes = append(recoveryCodes, code)
		}

		if err := store.Delete2faSetupSession(ctx, currentUser.Id); err != nil {
			return httpcomm.NewInternalServerError(util.Wrap(err), "Error while deleting setup session")
		}
		return nil
	})
	if err != nil {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Error while completing two factor authentication setup.")
	}

	return recoveryCodes, nil
}

func (uc *usersController) Disable2fa(ctx context.Context, currentUser *models.User, disable2faRequest *payloads.Disable2faRequest) error {
	if err := currentUser.ComparePasswords(disable2faRequest.Password); err != nil {
		return httpcomm.NewForbiddenError(util.Wrap(err), "Invalid credentials.")
	}

	err := uc.txnManager.WithTx(ctx, func(tx db.Querier) error {
		store := uc.usersStore.WithQuerier(tx)
		if err := store.DisableTwoFactorAuth(ctx, currentUser.Id); err != nil {
			return httpcomm.NewInternalServerError(util.Wrap(err), "Error while disabling two factor auth")
		}
		if err := store.DeleteRecoveryCodes(ctx, currentUser.Id); err != nil {
			return httpcomm.NewInternalServerError(util.Wrap(err), "Error while deleting recovery codes")
		}
		return nil
	})
	if err != nil {
		return httpcomm.NewInternalServerError(util.Wrap(err), "Error while disabling two factor authentication.")
	}

	return nil
}

func (uc *usersController) UpdatePassword(ctx context.Context, currentUser *models.User, updatePasswordRequest *payloads.UpdatePasswordRequest) error {
	if err := currentUser.ComparePasswords(updatePasswordRequest.CurrentPassword); err != nil {
		return httpcomm.NewForbiddenError(util.Wrap(err), "Invalid credentials.")
	}

	currentUser.Password = updatePasswordRequest.NewPassword
	if err := currentUser.HashPassword(); err != nil {
		return httpcomm.NewInternalServerError(util.Wrap(err), "Error while hashing password.")
	}

	if err := uc.usersStore.UpdatePassword(ctx, currentUser.Password, currentUser.Id); err != nil {
		return httpcomm.NewInternalServerError(util.Wrap(err), "Error while updating password.")
	}

	return nil
}

func (uc *usersController) RegenerateRecoveryCodes(ctx context.Context, currentUser *models.User) ([]*models.RecoveryCode, error) {
	recoveryCodes := make([]*models.RecoveryCode, 0, 10)
	err := uc.txnManager.WithTx(ctx, func(tx db.Querier) error {
		store := uc.usersStore.WithQuerier(tx)
		if err := store.DeleteRecoveryCodes(ctx, currentUser.Id); err != nil {
			return httpcomm.NewInternalServerError(util.Wrap(err), "Error while deleting recovery codes")
		}

		recoveryCodes = recoveryCodes[:0]
		for i := 0; i < 10; i++ {
			code, err := store.GenerateRecoveryCode(ctx, currentUser.Id)
			if err != nil {
				return httpcomm.NewInternalServerError(util.Wrap(err), "Error while generating recovery code")
			}
			recoveryCodes = append(recoveryCodes, code)
		}
		return nil
	})
	if err != nil {
		return nil, httpcomm.NewInternalServerError(util.Wrap(err), "Error while regenerating recovery codes.")
	}

	return recoveryCodes, nil
}
