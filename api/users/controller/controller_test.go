package controller

import (
	"context"
	"example/dashboard/api/models"
	"example/dashboard/api/users/mocks"
	"example/dashboard/config"
	"example/dashboard/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestUsersController_CreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockLogger := util.NewLogger()
	cfg := &config.Config{}

	userController := NewUsersController(cfg, mockStore, mockLogger)
	user := &models.User{
		Password: "password",
		Email:    "test5@email.com",
	}

	userWithtoken := &models.UserWithToken{
		User:  user,
		Token: "token",
	}
	mockStore.EXPECT().GetUserByEmail(context.TODO(), user.Email).Return(nil, pgx.ErrNoRows)
	mockStore.EXPECT().CreateUser(context.TODO(), user, gomock.Any()).Return(userWithtoken, nil)

	createdUser, err := userController.CreateUser(context.TODO(), user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

}

func TestUsersController_Login(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockLogger := util.NewLogger()
	cfg := &config.Config{}

	userController := NewUsersController(cfg, mockStore, mockLogger)

	loginUser := &models.User{
		Password: "password",
		Email:    "test5@email.com",
	}
	existingUser := &models.User{
		Password: "password",
		Email:    "test5@email.com",
	}

	mockStore.EXPECT().GetUserByEmail(context.TODO(), loginUser.Email).Return(existingUser, nil)

	userController.Login(context.TODO(), loginUser)
}
