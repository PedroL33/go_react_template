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
	CreateLoginSession(ctx context.Context, userId int, tx db.DbConn) (int, error)
	GetLoginSessionById(ctx context.Context, loginSessionId int, tx db.DbConn) (*models.LoginSession, error)
	DeleteLoginSessionByUserId(ctx context.Context, userId int, tx db.DbConn) error
	GetRecoveryCodesByUserId(ctx context.Context, userId int, tx db.DbConn) ([]*models.RecoveryCode, error)
	RedeemRecoveryCode(ctx context.Context, id int, tx db.DbConn) error
	UpdatePassword(ctx context.Context, password string, userId int, tx db.DbConn) error
}
