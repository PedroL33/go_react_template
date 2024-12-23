package users

import (
	"context"
	"example/dashboard/api/models"
	"example/dashboard/api/users/payloads"
)

type Controller interface {
	CreateUser(ctx context.Context, user *payloads.CreateUserRequest) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Begin2faSetupSession(ctx context.Context, currentUser *models.User) (string, error)
	Complete2faSetup(ctx context.Context, complete2faSetupRequest *payloads.Complete2faSetupRequest, currentUser *models.User) ([]*models.RecoveryCode, error)
	Disable2fa(ctx context.Context, currentUser *models.User) error
}
