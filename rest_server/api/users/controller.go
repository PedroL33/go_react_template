package users

import (
	"context"
	"example/template/rest_server/api/models"
	"example/template/rest_server/api/users/payloads"
)

//go:generate mockgen -source=controller.go -destination=./mocks/mock_controller.go -package=mocks

type Controller interface {
	CreateUser(ctx context.Context, user *payloads.CreateUserRequest) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*payloads.LoginResponse, error)
	Begin2faSetupSession(ctx context.Context, currentUser *models.User) (string, error)
	Complete2faSetup(ctx context.Context, complete2faSetupRequest *payloads.Complete2faSetupRequest, currentUser *models.User) ([]*models.RecoveryCode, error)
	Disable2fa(ctx context.Context, currentUser *models.User, disable2faRequest *payloads.Disable2faRequest) error
	VerifyLogin(ctx context.Context, verifyLoginRequest *payloads.LoginWithOptCodeRequest) (*payloads.LoginResponse, error)
	VerifyLoginWithRecoveryCode(ctx context.Context, verifyLoginRequest *payloads.LoginWithRecoveryCodeRequest) (*payloads.LoginResponse, error)
	UpdatePassword(ctx context.Context, currentUser *models.User, updatePasswordRequest *payloads.UpdatePasswordRequest) error
	RegenerateRecoveryCodes(ctx context.Context, currentUser *models.User) ([]*models.RecoveryCode, error)
}
