package users

import (
	"context"
	"example/dashboard/api/models"
)

type Controller interface {
	CreateUser(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
}
