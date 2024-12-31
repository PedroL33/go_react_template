package controller

import (
	"context"
	"database/sql"
	"example/dashboard/api/models"
	"example/dashboard/api/users/mocks"
	"example/dashboard/api/users/payloads"
	"example/dashboard/config"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUsersController_CreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockDbConn := mocks.NewMockDbConn(ctrl)
	cfg := &config.Config{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)
	createUserRequest := &payloads.CreateUserRequest{
		Password: "password",
		Email:    "test5@email.com",
	}

	user := &models.User{
		Email:    "test5@email.com",
		Password: "password",
	}

	mockStore.EXPECT().GetUserByEmail(context.TODO(), "test5@email.com", nil).Return(nil, pgx.ErrNoRows)
	mockStore.EXPECT().CreateUser(context.TODO(), gomock.Any(), mockDbConn).Return(user, nil)
	mockTxnManager.EXPECT().Begin(context.TODO()).Return(mockDbConn, nil)
	mockDbConn.EXPECT().Commit(context.TODO()).Return(nil)
	mockDbConn.EXPECT().Rollback(context.TODO()).Return(nil)
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
	cfg := &config.Config{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	existingUser := &models.User{
		Password: "password",
		Email:    "test5@email.com",
	}
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(existingUser.Password), bcrypt.DefaultCost)
	require.NoError(t, err)
	loginUser := &models.User{
		Password: string(hashedPw),
		Email:    "test5@email.com",
	}

	mockStore.EXPECT().GetUserByEmail(context.TODO(), gomock.Any(), nil).Return(loginUser, nil)

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
	cfg := &config.Config{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	twoFaSetupSession := &models.TwoFactorSetupSession{}
	twoFaSetupSession.PopulateSecretStringAndReturnBase64QrCode("test@email.com")
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

	mockStore.EXPECT().GetLoginSessionById(gomock.Any(), gomock.Any(), gomock.Any()).Return(loginSession, nil)
	mockStore.EXPECT().GetUserById(gomock.Any(), gomock.Any(), gomock.Any()).Return(currentUser, nil)

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
	cfg := &config.Config{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	verifyLoginRequest := &payloads.LoginWithRecoveryCodeRequest{
		RecoveryCode:   "test code",
		LoginSessionId: 1,
	}

	loginSession := &models.LoginSession{
		Expiration: time.Now().Add(1 * time.Minute),
	}

	var currentUser = &models.User{Id: 1}

	recoveryCodes := make([]*models.RecoveryCode, 0, 1)
	recoveryCodes = append(recoveryCodes, &models.RecoveryCode{Code: "test code"})

	mockStore.EXPECT().GetLoginSessionById(gomock.Any(), gomock.Any(), gomock.Any()).Return(loginSession, nil)
	mockStore.EXPECT().GetUserById(gomock.Any(), gomock.Any(), gomock.Any()).Return(currentUser, nil)
	mockStore.EXPECT().GetRecoveryCodesByUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(recoveryCodes, nil)
	mockStore.EXPECT().RedeemRecoveryCode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

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
	mockDbConn := mocks.NewMockDbConn(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	currentUser := &models.User{
		Email:    "test@email.com",
		Password: "test",
	}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	twoFASession := &models.TwoFactorSetupSession{}
	mockStore.EXPECT().Create2faSetupSession(context.TODO(), gomock.Any(), gomock.Any()).Return(twoFASession, nil)
	mockTxnManager.EXPECT().Begin(context.TODO()).Return(mockDbConn, nil)
	mockStore.EXPECT().Delete2faSetupSession(context.TODO(), gomock.Any(), gomock.Any()).Return(nil)
	mockDbConn.EXPECT().Commit(context.TODO()).Return(nil)
	mockDbConn.EXPECT().Rollback(context.TODO()).Return(nil)

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
	mockDbConn := mocks.NewMockDbConn(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	currentUser := &models.User{
		Email:    "test@email.com",
		Password: "password",
	}

	twoFaSetupSession := &models.TwoFactorSetupSession{}
	twoFaSetupSession.PopulateSecretStringAndReturnBase64QrCode(currentUser.Email)
	code, err := totp.GenerateCode(twoFaSetupSession.SecretString, time.Now())
	require.NoError(t, err)

	request := &payloads.Complete2faSetupRequest{
		OtpCode: code,
	}

	mockTxnManager.EXPECT().Begin(context.TODO()).Return(mockDbConn, nil)
	mockStore.EXPECT().Get2faSetupSessionByUserId(context.TODO(), gomock.Any(), gomock.Any()).Return(twoFaSetupSession, nil)
	mockDbConn.EXPECT().Rollback(context.TODO()).Return(nil)
	mockStore.EXPECT().EnableTwoFactorAuth(context.TODO(), gomock.Any(), gomock.Any()).Return(nil)
	mockStore.EXPECT().GenerateRecoveryCode(context.TODO(), gomock.Any(), gomock.Any()).Return(&models.RecoveryCode{}, nil).Times(10)
	mockStore.EXPECT().Delete2faSetupSession(context.TODO(), gomock.Any(), gomock.Any()).Return(nil)
	mockDbConn.EXPECT().Commit(context.TODO()).Return(nil)

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
	mockDbConn := mocks.NewMockDbConn(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)
	user := &models.User{
		Email:    "test@email.com",
		Password: "password",
	}

	disable2faRequest := &payloads.Disable2faRequest{Password: user.Password}

	err := user.HashPassword()
	require.NoError(t, err)

	mockTxnManager.EXPECT().Begin(context.TODO()).Return(mockDbConn, nil)
	mockStore.EXPECT().DisableTwoFactorAuth(context.TODO(), gomock.Any(), gomock.Any()).Return(nil)
	mockDbConn.EXPECT().Rollback(context.TODO()).Return(nil)
	mockStore.EXPECT().DeleteRecoveryCodes(context.TODO(), gomock.Any(), gomock.Any()).Return(nil)
	mockDbConn.EXPECT().Commit(context.TODO()).Return(nil)

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
	cfg := &config.Config{}

	currentUser := &models.User{Password: "test"}
	request := &payloads.UpdatePasswordRequest{CurrentPassword: currentUser.Password}

	err := currentUser.HashPassword()
	require.NoError(t, err)

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	mockStore.EXPECT().UpdatePassword(context.TODO(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err = userController.UpdatePassword(context.TODO(), currentUser, request)
	require.NoError(t, err)
}

func TestUsersController_RegenerateRecoveryCodes(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockTxnManager := mocks.NewMockTransactionManager(ctrl)
	mockDbConn := mocks.NewMockDbConn(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userController := NewUsersController(cfg, mockStore, mockTxnManager, mockLogger)

	mockTxnManager.EXPECT().Begin(context.TODO()).Return(mockDbConn, nil)
	mockStore.EXPECT().DeleteRecoveryCodes(context.TODO(), gomock.Any(), gomock.Any()).Return(nil)
	mockStore.EXPECT().GenerateRecoveryCode(context.TODO(), gomock.Any(), gomock.Any()).Return(&models.RecoveryCode{}, nil).Times(10)
	mockDbConn.EXPECT().Commit(context.TODO()).Return(nil)
	mockDbConn.EXPECT().Rollback((context.TODO())).Return(nil)

	currentUser := &models.User{Password: "test"}
	codes, err := userController.RegenerateRecoveryCodes(context.TODO(), currentUser)
	require.NoError(t, err)
	require.NotNil(t, codes)
}
