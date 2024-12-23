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
		Email:    "test@email.com",
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
	response := &models.UserWithToken{}
	err = json.NewDecoder(rec.Body).Decode(response)

	require.NoError(t, err)
}

func TestUsersHandlers_Login(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	user := &models.User{
		Email:    "test@email.com",
		Password: "testpassword",
	}

	userWithToken := &models.UserWithToken{
		User: &models.User{
			Id: 1,
		},
	}

	body, err := json.Marshal(user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockLogger.EXPECT().HttpSuccess(req)
	mockController.EXPECT().Login(gomock.Any(), user).Return(userWithToken, nil)

	userHandler.Login(rec, req)

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
		Email: "test@email.com",
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
		Email: "test@email.com",
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
		Email: "test@email.com",
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

	mockController.EXPECT().Disable2fa(ctx, user).Return(nil)
	mockLogger.EXPECT().HttpSuccess(req)

	userHandler.Disable2fa(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}
