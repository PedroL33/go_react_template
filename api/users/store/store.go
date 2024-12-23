package store

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"example/dashboard/api/base"
	"example/dashboard/api/db"
	"example/dashboard/api/models"
	"example/dashboard/api/users"
	"example/dashboard/util"

	"github.com/pkg/errors"
)

type userStore struct {
	db db.DbConn
	*base.Store
}

func NewUsersStore(db db.DbConn) users.Store {
	return &userStore{
		db:    db,
		Store: &base.Store{DB: db},
	}
}

func (s *userStore) CreateUser(ctx context.Context, user *models.User, tx db.DbConn) (*models.User, error) {
	db := db.GetDb(s.db, tx)

	row := db.QueryRow(
		ctx,
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *",
		user.Email,
		user.Password,
	)

	createdUser := &models.User{}

	if err := util.ScanRowIntoStruct(row, createdUser); err != nil {
		return nil, errors.Wrap(err, "userStore.CreateUser")
	}

	return createdUser, nil
}

func (s *userStore) GetUserByEmail(ctx context.Context, email string, tx db.DbConn) (*models.User, error) {
	db := db.GetDb(s.db, tx)

	row := db.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email)
	foundUser := &models.User{}
	row.Scan(&foundUser.Id, &foundUser.Email, &foundUser.Password, &foundUser.IsTwoFactorEnabled, &foundUser.TwoFactorSecret, &foundUser.CreatedAt, &foundUser.UpdatedAt)
	// if err := util.ScanRowIntoStruct(row, foundUser); err != nil {
	// 	return nil, errors.Wrap(err, "userStore.GetUserByEmail")
	// }

	return foundUser, nil
}

func (s *userStore) GetUserById(ctx context.Context, userId int, tx db.DbConn) (*models.User, error) {
	db := db.GetDb(s.db, tx)

	row := db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", userId)

	foundUser := &models.User{}
	if err := util.ScanRowIntoStruct(row, foundUser); err != nil {
		return nil, errors.Wrap(err, "userStore.GetUserById")
	}

	return foundUser, nil
}

func (s *userStore) Create2faSetupSession(ctx context.Context, session *models.TwoFactorSetupSession, tx db.DbConn) (*models.TwoFactorSetupSession, error) {
	db := db.GetDb(s.db, tx)

	row := db.QueryRow(
		ctx,
		"INSERT INTO two_factor_setup_sessions (user_id, secret_string, expiration_timestamp) VALUES ($1, $2, $3) RETURNING *",
		session.UserId,
		session.SecretString,
		session.Expiration,
	)

	createdSession := &models.TwoFactorSetupSession{}
	if err := util.ScanRowIntoStruct(row, createdSession); err != nil {
		return nil, errors.Wrap(err, "userStore.Create2faSetupSession")
	}

	return createdSession, nil
}

func (s *userStore) Get2faSetupSessionByUserId(ctx context.Context, userId int, tx db.DbConn) (*models.TwoFactorSetupSession, error) {
	db := db.GetDb(s.db, tx)

	row := db.QueryRow(ctx, "SELECT * FROM two_factor_setup_sessions WHERE user_id = $1", userId)

	setupSession := &models.TwoFactorSetupSession{}
	if err := util.ScanRowIntoStruct(row, setupSession); err != nil {
		return nil, errors.Wrap(err, "userStore.Get2faSetupSessionByUserId")
	}

	return setupSession, nil
}

func (s *userStore) EnableTwoFactorAuth(ctx context.Context, setupSession *models.TwoFactorSetupSession, tx db.DbConn) error {
	db := db.GetDb(s.db, tx)
	if _, err := db.Exec(
		ctx,
		"UPDATE users SET two_factor_secret = $1, is_two_factor_enabled = $2 WHERE id=$3",
		setupSession.SecretString,
		true,
		setupSession.UserId,
	); err != nil {
		return errors.Wrap(err, "userStore.EnableTwoFactorAuth")
	}

	return nil
}

func (s *userStore) DisableTwoFactorAuth(ctx context.Context, userId int, tx db.DbConn) error {
	db := db.GetDb(s.db, tx)
	if _, err := db.Exec(
		ctx,
		"UPDATE users SET two_factor_secret = $1, is_two_factor_enabled = $2 WHERE id=$3",
		nil,
		false,
		userId,
	); err != nil {
		return errors.Wrap(err, "userStore.DisableTwoFactorAuth")
	}

	return nil
}

func (s *userStore) GenerateRecoveryCode(ctx context.Context, userId int, tx db.DbConn) (*models.RecoveryCode, error) {
	db := db.GetDb(s.db, tx)

	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "userStore.DisableTwoFactorAuth")
	}

	code := base64.URLEncoding.EncodeToString(bytes)[:8]

	row := db.QueryRow(ctx, "INSERT INTO recovery_codes (user_id, code) VALUES ($1, $2) RETURNING *", userId, code)

	recoveryCode := &models.RecoveryCode{}
	if err := util.ScanRowIntoStruct(row, recoveryCode); err != nil {
		return nil, errors.Wrap(err, "userStore.GenerateRecoveryCodes")
	}

	return recoveryCode, nil
}

func (s *userStore) Delete2faSetupSession(ctx context.Context, userId int, tx db.DbConn) error {
	db := db.GetDb(s.db, tx)

	if _, err := db.Exec(
		ctx,
		"DELETE FROM two_factor_setup_sessions WHERE user_id = $1",
		userId,
	); err != nil {
		return errors.Wrap(err, "userStore.Delete2faSetupSession")
	}

	return nil
}

func (s *userStore) DeleteRecoveryCodes(ctx context.Context, userId int, tx db.DbConn) error {
	db := db.GetDb(s.db, tx)

	if _, err := db.Exec(
		ctx,
		"DELETE FROM recovery_codes WHERE user_id = $1",
		userId,
	); err != nil {
		return errors.Wrap(err, "userStore.DeleteRecoveryCodes")
	}

	return nil
}
