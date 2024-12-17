package store

import (
	"context"
	"example/dashboard/api/models"
	"example/dashboard/api/users/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUsersStore_CreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbConn := mocks.NewMockDbConn(ctrl)
	userStore := NewUsersStore(mockDbConn)
	tx := mocks.NewMockDbConn(ctrl)
	user := &models.User{
		Email:    "",
		Password: "",
	}
	row := mocks.NewMockDbRow(ctrl)

	mockDbConn.EXPECT().Begin(context.TODO()).Return(tx, nil)
	tx.EXPECT().QueryRow(context.TODO(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
	tx.EXPECT().Rollback(context.TODO()).Return(nil)
	tx.EXPECT().Commit(context.TODO()).Return(nil)
	row.EXPECT().Scan(gomock.Any()).Return(nil)

	createdUser, err := userStore.CreateUser(context.TODO(), user, func() (string, error) { return "", nil })

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

	user, err := userStore.GetUserByEmail(context.TODO(), "test@email.com")

	require.NoError(t, err)
	require.NotNil(t, user)
}
