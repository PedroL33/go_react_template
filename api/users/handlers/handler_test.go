package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"example/dashboard/api/middleware"
	"example/dashboard/api/models"
	"example/dashboard/api/users/mocks"
	"example/dashboard/api/users/payloads"
	"example/dashboard/config"
	"example/dashboard/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUsersHandlers_Create(t *testing.T) {

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	user := &models.User{
		Username: "username",
		Password: "testpassword",
	}

	userWithToken := &models.UserWithToken{
		User: &models.User{
			Id: 1,
		},
	}

	body, err := json.Marshal(user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockController.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(userWithToken, nil)
	mockLogger.EXPECT().HttpSuccess(req)
	userHandler.Create(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUsersHandlers_Login(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	userLoginRequest := &payloads.UserLoginRequest{}

	userWithToken := &payloads.LoginResponse{}
	body, err := json.Marshal(userLoginRequest)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockLogger.EXPECT().HttpSuccess(req)
	mockController.EXPECT().Login(gomock.Any(), gomock.Any()).Return(userWithToken, nil)

	userHandler.Login(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUsersHandlers_VerifyLogin(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	userLoginRequest := &payloads.UserLoginRequest{}

	loginResponse := &payloads.LoginResponse{}
	body, err := json.Marshal(userLoginRequest)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login_otp", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockLogger.EXPECT().HttpSuccess(req)
	mockController.EXPECT().VerifyLogin(gomock.Any(), gomock.Any()).Return(loginResponse, nil)

	userHandler.VerifyLogin(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUsersHandlers_VerifyLoginWithRecoveryCode(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	userLoginRequest := &payloads.UserLoginRequest{}

	loginResponse := &payloads.LoginResponse{}
	body, err := json.Marshal(userLoginRequest)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login_recovery_code", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockLogger.EXPECT().HttpSuccess(req)
	mockController.EXPECT().VerifyLoginWithRecoveryCode(gomock.Any(), gomock.Any()).Return(loginResponse, nil)

	userHandler.VerifyLoginWithRecoveryCode(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUsersHandlers_Begin2faSetup(t *testing.T) {
	var err error
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	user := &models.User{
		Username: "username",
	}
	var token string
	token, err = util.CreateToken(cfg, user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/begin2fa", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer: "+token)
	rec := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middleware.CurrentUserKey, user)
	req = req.WithContext(ctx)

	mockController.EXPECT().Begin2faSetupSession(ctx, gomock.Any()).Return("test", nil)
	mockLogger.EXPECT().HttpSuccess(req)

	userHandler.Begin2faSetup(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUsersHandlers_Complete2faSetup(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	user := &models.User{
		Username: "username",
	}

	request := &payloads.Complete2faSetupRequest{
		OtpCode: "test",
	}

	body, err := json.Marshal(request)
	require.NoError(t, err)

	token, err := util.CreateToken(cfg, user)
	require.NoError(t, err)

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)
	req := httptest.NewRequest(http.MethodPost, "/complete2fa", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer: "+token)
	rec := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middleware.CurrentUserKey, user)
	req = req.WithContext(ctx)

	retVal := []*models.RecoveryCode{}
	mockController.EXPECT().Complete2faSetup(ctx, request, user).Return(retVal, nil)
	mockLogger.EXPECT().HttpSuccess(req)

	userHandler.Complete2faSetup(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUsersHandlers_Disable2fa(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	user := &models.User{
		Username: "username",
	}

	request := &payloads.Complete2faSetupRequest{
		OtpCode: "test",
	}

	body, err := json.Marshal(request)
	require.NoError(t, err)

	token, err := util.CreateToken(cfg, user)
	require.NoError(t, err)

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)
	req := httptest.NewRequest(http.MethodPost, "/disable2fa", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer: "+token)
	rec := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middleware.CurrentUserKey, user)
	req = req.WithContext(ctx)

	mockController.EXPECT().Disable2fa(ctx, user, &payloads.Disable2faRequest{}).Return(nil)
	mockLogger.EXPECT().HttpSuccess(req)

	userHandler.Disable2fa(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUsersHandlers_ChangePassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	user := &models.User{
		Username: "username",
	}

	request := &payloads.UpdatePasswordRequest{
		NewPassword:     "testtest1",
		CurrentPassword: "testtest",
	}

	body, err := json.Marshal(request)
	require.NoError(t, err)
	token, err := util.CreateToken(cfg, user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/change_password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer: "+token)
	rec := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middleware.CurrentUserKey, user)
	req = req.WithContext(ctx)

	mockController.EXPECT().UpdatePassword(ctx, user, request).Return(nil)
	mockLogger.EXPECT().HttpSuccess(req)

	userHandler.ChangePassword(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUsersHandlers_RegenerateRecoveryCodes(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	user := &models.User{
		Username: "username",
	}

	token, err := util.CreateToken(cfg, user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/change_password", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer: "+token)
	rec := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middleware.CurrentUserKey, user)
	req = req.WithContext(ctx)

	mockController.EXPECT().RegenerateRecoveryCodes(ctx, user).Return([]*models.RecoveryCode{}, nil)
	mockLogger.EXPECT().HttpSuccess(req)

	userHandler.RegenerateRecoveryCodes(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}
