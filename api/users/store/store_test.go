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
		Username: "",
		Password: "",
	}
	row := mocks.NewMockDbRow(ctrl)

	mockDbConn.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	createdUser, err := userStore.CreateUser(context.TODO(), user, nil)

	require.NoError(t, err)
	require.NotNil(t, createdUser)
}

func TestUsersStore_GetUserByUsername(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)
	row := mocks.NewMockDbRow(ctrl)

	mockDbconn.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	user, err := userStore.GetUserByUsername(context.TODO(), "username", nil)

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

func TestUsersStore_GetRecoveryCodesByUserId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	rows := mocks.NewMockDbRows(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().Query(context.TODO(), gomock.Any(), gomock.Any()).Return(rows, nil)
	rows.EXPECT().Next().Return(true).Times(9)
	rows.EXPECT().Next().Return(false).Times(1)
	rows.EXPECT().Scan(gomock.Any()).Return(nil).Times(9)

	codes, err := userStore.GetRecoveryCodesByUserId(context.TODO(), 1, mockDbconn)
	require.NoError(t, err)
	require.NotNil(t, codes)
}

func TestUsersStore_RedeemRecoveryCode(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.RedeemRecoveryCode(context.TODO(), 1, mockDbconn)
	require.NoError(t, err)
}

func TestUsersStore_CreateLoginSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	row := mocks.NewMockDbRow(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	session, err := userStore.CreateLoginSession(context.TODO(), 1, mockDbconn)
	require.NoError(t, err)
	require.NotNil(t, session)
}

func TestUsersStore_GetLoginSessionById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)
	row := mocks.NewMockDbRow(ctrl)

	mockDbconn.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	session, err := userStore.GetLoginSessionById(context.TODO(), 1, mockDbconn)
	require.NoError(t, err)
	require.NotNil(t, session)
}

func TestUsersStore_DeleteLoginSessionByUserId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.DeleteLoginSessionByUserId(context.TODO(), 1, mockDbconn)
	require.NoError(t, err)
}

func TestUsersStore_UpdatePassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbconn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbconn)

	mockDbconn.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.UpdatePassword(context.TODO(), "test", 1, mockDbconn)
	require.NoError(t, err)
}
