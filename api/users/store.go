package users

import (
	"context"
	"example/dashboard/api/db"
	"example/dashboard/api/models"
)

//go:generate mockgen -source=store.go -destination=./mocks/mock_store.go -package=mocks

// Store runs user-related queries against whatever Querier it is bound to
// (the pool, a pinned connection, or a transaction). Callers rebind the store
// with WithQuerier when they need tx/conn-scoped execution.
type Store interface {
	// WithQuerier returns a new Store view bound to q. The original store is
	// unchanged. Typically used inside a TransactionManager.WithTx callback
	// or after Pool.Acquire.
	WithQuerier(q db.Querier) Store

	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	Create2faSetupSession(ctx context.Context, session *models.TwoFactorSetupSession) (*models.TwoFactorSetupSession, error)
	GetUserById(ctx context.Context, userId int) (*models.User, error)
	Get2faSetupSessionByUserId(ctx context.Context, userId int) (*models.TwoFactorSetupSession, error)
	EnableTwoFactorAuth(ctx context.Context, setupSession *models.TwoFactorSetupSession) error
	DisableTwoFactorAuth(ctx context.Context, userId int) error
	GenerateRecoveryCode(ctx context.Context, userId int) (*models.RecoveryCode, error)
	Delete2faSetupSession(ctx context.Context, userId int) error
	DeleteRecoveryCodes(ctx context.Context, userId int) error
	CreateLoginSession(ctx context.Context, userId int) (int, error)
	GetLoginSessionById(ctx context.Context, loginSessionId int) (*models.LoginSession, error)
	DeleteLoginSessionByUserId(ctx context.Context, userId int) error
	GetRecoveryCodesByUserId(ctx context.Context, userId int) ([]*models.RecoveryCode, error)
	RedeemRecoveryCode(ctx context.Context, id int) error
	UpdatePassword(ctx context.Context, password string, userId int) error
}
