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

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)

	user := &models.User{
		Username: "",
		Password: "",
	}
	row := mocks.NewMockDbRow(ctrl)

	mockQuerier.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	createdUser, err := userStore.CreateUser(context.TODO(), user)

	require.NoError(t, err)
	require.NotNil(t, createdUser)
}

func TestUsersStore_GetUserByUsername(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)
	row := mocks.NewMockDbRow(ctrl)

	mockQuerier.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	user, err := userStore.GetUserByUsername(context.TODO(), "username")

	require.NoError(t, err)
	require.NotNil(t, user)
}

func TestUsersStore_Create2faSetupSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)
	row := mocks.NewMockDbRow(ctrl)

	mockQuerier.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	setupSession, err := userStore.Create2faSetupSession(context.TODO(), &models.TwoFactorSetupSession{})

	require.NoError(t, err)
	require.NotNil(t, setupSession)
}

func TestUsersStore_Get2faSetupSessionByUserId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)
	row := mocks.NewMockDbRow(ctrl)

	mockQuerier.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	setupSession, err := userStore.Get2faSetupSessionByUserId(context.TODO(), 1)
	require.NoError(t, err)
	require.NotNil(t, setupSession)
}

func TestUsersStore_EnableTwoFactorAuth(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.EnableTwoFactorAuth(context.TODO(), &models.TwoFactorSetupSession{})
	require.NoError(t, err)
}

func TestUsersStore_DisableTwoFactorAuth(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.DisableTwoFactorAuth(context.TODO(), 1)
	require.NoError(t, err)
}

func TestUsersStore_GenerateRecoveryCode(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	row := mocks.NewMockDbRow(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	recoveryCodes, err := userStore.GenerateRecoveryCode(context.TODO(), 1)
	require.NoError(t, err)
	require.NotNil(t, recoveryCodes)
}

func TestUsersStore_Delete2faSetupSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.Delete2faSetupSession(context.TODO(), 1)
	require.NoError(t, err)
}

func TestUsersStore_DeleteRecoveryCodes(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.DeleteRecoveryCodes(context.TODO(), 1)
	require.NoError(t, err)
}

func TestUsersStore_GetRecoveryCodesByUserId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	rows := mocks.NewMockDbRows(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().Query(context.TODO(), gomock.Any(), gomock.Any()).Return(rows, nil)
	rows.EXPECT().Next().Return(true).Times(9)
	rows.EXPECT().Next().Return(false).Times(1)
	rows.EXPECT().Scan(gomock.Any()).Return(nil).Times(9)

	codes, err := userStore.GetRecoveryCodesByUserId(context.TODO(), 1)
	require.NoError(t, err)
	require.NotNil(t, codes)
}

func TestUsersStore_RedeemRecoveryCode(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.RedeemRecoveryCode(context.TODO(), 1)
	require.NoError(t, err)
}

func TestUsersStore_CreateLoginSession(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	row := mocks.NewMockDbRow(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	session, err := userStore.CreateLoginSession(context.TODO(), 1)
	require.NoError(t, err)
	require.NotNil(t, session)
}

func TestUsersStore_GetLoginSessionById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)
	row := mocks.NewMockDbRow(ctrl)

	mockQuerier.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	session, err := userStore.GetLoginSessionById(context.TODO(), 1)
	require.NoError(t, err)
	require.NotNil(t, session)
}

func TestUsersStore_DeleteLoginSessionByUserId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.DeleteLoginSessionByUserId(context.TODO(), 1)
	require.NoError(t, err)
}

func TestUsersStore_UpdatePassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockQuerier(ctrl)
	userStore := NewUsersStore(mockQuerier)

	mockQuerier.EXPECT().Exec(context.TODO(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	err := userStore.UpdatePassword(context.TODO(), "test", 1)
	require.NoError(t, err)
}

func TestUsersStore_WithQuerier_ReturnsIndependentView(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	poolMock := mocks.NewMockQuerier(ctrl)
	txMock := mocks.NewMockQuerier(ctrl)
	row := mocks.NewMockDbRow(ctrl)

	baseStore := NewUsersStore(poolMock)
	txStore := baseStore.WithQuerier(txMock)

	// The re-bound store should hit the tx, not the pool.
	txMock.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	user, err := txStore.GetUserByUsername(context.TODO(), "foo")
	require.NoError(t, err)
	require.NotNil(t, user)

	// Base store still points at the pool (unchanged).
	poolMock.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).Return(nil)
	_, err = baseStore.GetUserByUsername(context.TODO(), "foo")
	require.NoError(t, err)
}
