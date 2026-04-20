package controller

import (
	"context"
	"database/sql"
	"example/template/rest_server/api/db"
	"example/template/rest_server/api/models"
	"example/template/rest_server/api/users/mocks"
	"example/template/rest_server/api/users/payloads"
	"example/template/rest_server/config"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// expectTxWithStore wires up the common pattern used by controller methods
// that run work inside a transaction:
//  1. WithTx invokes its callback with the provided querier.
//  2. store.WithQuerier(tx) returns the same mock store, so a single set of
//     EXPECT() calls on mockStore covers both pool- and tx-scoped methods.
func expectTxWithStore(mockTxnManager *mocks.MockTransactionManager, mockStore *mocks.MockStore, mockQuerier *mocks.MockQuerier) {
	mockTxnManager.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, fn func(db.Querier) error) error {
			return fn(mockQuerier)
		},
	)
	mockStore.EXPECT().WithQuerier(mockQuerier).Return(mockStore)
}

func TestUsersController_CreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockQuerier := mocks.NewMockQuerier(ctrl)
	cfg := &config.AppConfig{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)
	createUserRequest := &payloads.CreateUserRequest{
		Password: "password",
		Username: "username",
	}

	user := &models.User{
		Username: "username",
		Password: "password",
	}

	mockStore.EXPECT().GetUserByUsername(context.TODO(), "username").Return(nil, pgx.ErrNoRows)
	expectTxWithStore(mockTxnManager, mockStore, mockQuerier)
	mockStore.EXPECT().CreateUser(context.TODO(), gomock.Any()).Return(user, nil)

	createdUser, err := userController.CreateUser(context.TODO(), createUserRequest)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
}

func TestUsersController_Login(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.AppConfig{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	existingUser := &models.User{
		Password: "password",
		Username: "username",
	}
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(existingUser.Password), bcrypt.DefaultCost)
	require.NoError(t, err)
	loginUser := &models.User{
		Password: string(hashedPw),
		Username: "username",
	}

	mockStore.EXPECT().GetUserByUsername(context.TODO(), gomock.Any()).Return(loginUser, nil)

	loggedInUser, err := userController.Login(context.TODO(), existingUser)
	require.NoError(t, err)
	require.NotNil(t, loggedInUser)
}

func TestUsersController_VerifyLogin(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.AppConfig{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	twoFaSetupSession := &models.TwoFactorSetupSession{}
	twoFaSetupSession.PopulateSecretStringAndReturnBase64QrCode("username")
	code, err := totp.GenerateCode(twoFaSetupSession.SecretString, time.Now())
	require.NoError(t, err)
	verifyLoginRequest := &payloads.LoginWithOptCodeRequest{
		OtpCode:        code,
		LoginSessionId: 1,
	}

	loginSession := &models.LoginSession{
		Expiration: time.Now().Add(1 * time.Minute),
	}

	var currentUser = &models.User{
		TwoFactorSecret: sql.NullString{
			Valid:  true,
			String: twoFaSetupSession.SecretString,
		},
	}

	mockStore.EXPECT().GetLoginSessionById(gomock.Any(), gomock.Any()).Return(loginSession, nil)
	mockStore.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(currentUser, nil)

	loginResponse, err := userController.VerifyLogin(context.TODO(), verifyLoginRequest)
	require.NoError(t, err)
	require.NotNil(t, loginResponse)
}

func TestUsersController_VerifyLoginWithRecoveryCode(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.AppConfig{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	verifyLoginRequest := &payloads.LoginWithRecoveryCodeRequest{
		RecoveryCode:   "test code",
		LoginSessionId: 1,
	}

	loginSession := &models.LoginSession{
		Expiration: time.Now().Add(1 * time.Minute),
	}

	var currentUser = &models.User{Id: 1}

	recoveryCodes := []*models.RecoveryCode{{Code: "test code"}}

	mockStore.EXPECT().GetLoginSessionById(gomock.Any(), gomock.Any()).Return(loginSession, nil)
	mockStore.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(currentUser, nil)
	mockStore.EXPECT().GetRecoveryCodesByUserId(gomock.Any(), gomock.Any()).Return(recoveryCodes, nil)
	mockStore.EXPECT().RedeemRecoveryCode(gomock.Any(), gomock.Any()).Return(nil)

	loginResponse, err := userController.VerifyLoginWithRecoveryCode(context.TODO(), verifyLoginRequest)
	require.NoError(t, err)
	require.NotNil(t, loginResponse)
}

func TestUsersController_Begin2faSetupSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockQuerier := mocks.NewMockQuerier(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.AppConfig{}

	currentUser := &models.User{
		Username: "username",
		Password: "test",
	}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	twoFASession := &models.TwoFactorSetupSession{}
	expectTxWithStore(mockTxnManager, mockStore, mockQuerier)
	mockStore.EXPECT().Delete2faSetupSession(context.TODO(), gomock.Any()).Return(nil)
	mockStore.EXPECT().Create2faSetupSession(context.TODO(), gomock.Any()).Return(twoFASession, nil)

	code, err := userController.Begin2faSetupSession(context.TODO(), currentUser)
	require.NoError(t, err)
	require.NotNil(t, code)
}

func TestUsersController_Complete2faSetup(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockQuerier := mocks.NewMockQuerier(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.AppConfig{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	currentUser := &models.User{
		Username: "username",
		Password: "password",
	}

	twoFaSetupSession := &models.TwoFactorSetupSession{}
	twoFaSetupSession.PopulateSecretStringAndReturnBase64QrCode(currentUser.Username)
	code, err := totp.GenerateCode(twoFaSetupSession.SecretString, time.Now())
	require.NoError(t, err)

	request := &payloads.Complete2faSetupRequest{
		OtpCode: code,
	}

	mockStore.EXPECT().Get2faSetupSessionByUserId(context.TODO(), gomock.Any()).Return(twoFaSetupSession, nil)
	expectTxWithStore(mockTxnManager, mockStore, mockQuerier)
	mockStore.EXPECT().EnableTwoFactorAuth(context.TODO(), gomock.Any()).Return(nil)
	mockStore.EXPECT().GenerateRecoveryCode(context.TODO(), gomock.Any()).Return(&models.RecoveryCode{}, nil).Times(10)
	mockStore.EXPECT().Delete2faSetupSession(context.TODO(), gomock.Any()).Return(nil)

	codes, err := userController.Complete2faSetup(context.TODO(), request, currentUser)
	require.NoError(t, err)
	require.NotNil(t, codes)
}

func TestUsersController_Disable2fa(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockQuerier := mocks.NewMockQuerier(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.AppConfig{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)
	user := &models.User{
		Username: "username",
		Password: "password",
	}

	disable2faRequest := &payloads.Disable2faRequest{Password: user.Password}

	err := user.HashPassword()
	require.NoError(t, err)

	expectTxWithStore(mockTxnManager, mockStore, mockQuerier)
	mockStore.EXPECT().DisableTwoFactorAuth(context.TODO(), gomock.Any()).Return(nil)
	mockStore.EXPECT().DeleteRecoveryCodes(context.TODO(), gomock.Any()).Return(nil)

	err = userController.Disable2fa(context.TODO(), user, disable2faRequest)
	require.NoError(t, err)
}

func TestUsersController_UpdatePassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.AppConfig{}

	currentUser := &models.User{Password: "test"}
	request := &payloads.UpdatePasswordRequest{CurrentPassword: currentUser.Password}

	err := currentUser.HashPassword()
	require.NoError(t, err)

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	mockStore.EXPECT().UpdatePassword(context.TODO(), gomock.Any(), gomock.Any()).Return(nil)

	err = userController.UpdatePassword(context.TODO(), currentUser, request)
	require.NoError(t, err)
}

func TestUsersController_RegenerateRecoveryCodes(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockQuerier := mocks.NewMockQuerier(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.AppConfig{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	expectTxWithStore(mockTxnManager, mockStore, mockQuerier)
	mockStore.EXPECT().DeleteRecoveryCodes(context.TODO(), gomock.Any()).Return(nil)
	mockStore.EXPECT().GenerateRecoveryCode(context.TODO(), gomock.Any()).Return(&models.RecoveryCode{}, nil).Times(10)

	currentUser := &models.User{Password: "test"}
	codes, err := userController.RegenerateRecoveryCodes(context.TODO(), currentUser)
	require.NoError(t, err)
	require.NotNil(t, codes)
}
