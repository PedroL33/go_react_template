package handlers

import (
	"bytes"
	"encoding/json"
	"example/dashboard/api/models"
	"example/dashboard/api/users/mocks"
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
	mockLogger := util.NewLogger()
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	user := &models.User{
		Email:    "test@email.com",
		Password: "testpassword",
	}

	userWithToken := &models.UserWithToken{
		User: &models.User{
			Id: "1",
		},
	}

	body, err := json.Marshal(user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockController.EXPECT().CreateUser(gomock.Any(), user).Return(userWithToken, nil)

	userHandler.Create(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	responseUser := &models.UserWithToken{}
	err = json.NewDecoder(rec.Body).Decode(responseUser)
	require.NoError(t, err)
	require.Equal(t, userWithToken, responseUser)
}

func TestUsersHandlers_Login(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockController := mocks.NewMockController(ctrl)
	mockLogger := util.NewLogger()
	cfg := &config.Config{}

	userHandler := NewUsersHandlers(cfg, mockController, mockLogger)

	user := &models.User{
		Email:    "test@email.com",
		Password: "testpassword",
	}

	userWithToken := &models.UserWithToken{
		User: &models.User{
			Id: "1",
		},
	}

	body, err := json.Marshal(user)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockController.EXPECT().Login(gomock.Any(), user).Return(userWithToken, nil)

	userHandler.Login(rec, req)

	require.Equal(t, http.StatusAccepted, rec.Code)
	responseUser := &models.UserWithToken{}
	err = json.NewDecoder(rec.Body).Decode(responseUser)
	require.NoError(t, err)
	require.Equal(t, userWithToken, responseUser)
}
