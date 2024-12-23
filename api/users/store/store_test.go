package store

import (
	"context"
	"example/dashboard/api/models"
	"example/dashboard/api/users/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestUsersStore_CreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbConn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbConn)

	user := &models.User{
		Email:    "",
		Password: "",
	}
	row := mocks.NewMockDbRow(ctrl)

	mockDbConn.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	createdUser, err := userStore.CreateUser(context.TODO(), user, nil)

	require.NoError(t, err)
	require.NotNil(t, createdUser)
}

func TestUsersStore_GetUserById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)
	row := mocks.NewMockDbRow(ctrl)

	mockDbconn.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	user, err := userStore.GetUserByEmail(context.TODO(), "test@email.com", nil)

	require.NoError(t, err)
	require.NotNil(t, user)
}

func TestUsersStore_Create2faSetupSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)
	row := mocks.NewMockDbRow(ctrl)

	mockDbconn.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	setupSession, err := userStore.Create2faSetupSession(context.TODO(), &models.TwoFactorSetupSession{}, nil)

	require.NoError(t, err)
	require.NotNil(t, setupSession)
}

func TestUsersStore_Get2faSetupSessionByUserId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)
	row := mocks.NewMockDbRow(ctrl)

	mockDbconn.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	setupSession, err := userStore.Get2faSetupSessionByUserId(context.TODO(), 1, nil)
	require.NoError(t, err)
	require.NotNil(t, setupSession)
}

func TestUsersStore_EnableTwoFactorAuth(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.EnableTwoFactorAuth(context.TODO(), &models.TwoFactorSetupSession{}, mockDbconn)
	require.NoError(t, err)
}

func TestUsersStore_DisableTwoFactorAuth(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.DisableTwoFactorAuth(context.TODO(), 1, mockDbconn)
	require.NoError(t, err)
}

func TestUsersStore_GenerateRecoveryCode(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	row := mocks.NewMockDbRow(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	recoveryCodes, err := userStore.GenerateRecoveryCode(context.TODO(), 1, mockDbconn)
	require.NoError(t, err)
	require.NotNil(t, recoveryCodes)
}

func TestUsersStore_Delete2faSetupSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.Delete2faSetupSession(context.TODO(), 1, mockDbconn)
	require.NoError(t, err)
}

func TestUsersStore_DeleteRecoveryCodes(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.DeleteRecoveryCodes(context.TODO(), 1, mockDbconn)
	require.NoError(t, err)
}
