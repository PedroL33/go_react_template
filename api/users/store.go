package users

import (
	"context"
	"example/dashboard/api/models"
)

type Store interface {
	CreateUser(ctx context.Context, user *models.User, createToken func() (string, error)) (*models.UserWithToken, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
