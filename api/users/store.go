package users

import (
	"context"
	"example/dashboard/api/db"
	"example/dashboard/api/models"
)

type Store interface {
	CreateUser(ctx context.Context, user *models.User, tx db.DbConn) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string, tx db.DbConn) (*models.User, error)
	Create2faSetupSession(ctx context.Context, session *models.TwoFactorSetupSession, tx db.DbConn) (*models.TwoFactorSetupSession, error)
	GetUserById(ctx context.Context, userId int, tx db.DbConn) (*models.User, error)
	Get2faSetupSessionByUserId(ctx context.Context, userId int, tx db.DbConn) (*models.TwoFactorSetupSession, error)
	EnableTwoFactorAuth(ctx context.Context, setupSession *models.TwoFactorSetupSession, tx db.DbConn) error
	DisableTwoFactorAuth(ctx context.Context, userId int, tx db.DbConn) error
	GenerateRecoveryCode(ctx context.Context, userId int, tx db.DbConn) (*models.RecoveryCode, error)
	Delete2faSetupSession(ctx context.Context, userId int, tx db.DbConn) error
	DeleteRecoveryCodes(ctx context.Context, userId int, tx db.DbConn) error
}
