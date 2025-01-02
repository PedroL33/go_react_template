package store

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"example/dashboard/api/db"
	"example/dashboard/api/models"
	"example/dashboard/api/users"
	"example/dashboard/util"
	"time"

	"github.com/pkg/errors"
)

type userStore struct {
	db db.DbConn
}

func NewUsersStore(db db.DbConn) users.Store {
	return &userStore{
		db: db,
	}
}

func (s *userStore) CreateUser(ctx context.Context, user *models.User, tx db.DbConn) (*models.User, error) {
	db := db.GetDb(s.db, tx)
	row := db.QueryRow(
		ctx,
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *",
		user.Username,
		user.Password,
	)

	createdUser := &models.User{}

	if err := util.ScanRowIntoStruct(row, createdUser); err != nil {
		return nil, errors.Wrap(err, "userStore.CreateUser")
	}

	return createdUser, nil
}

func (s *userStore) GetUserByUsername(ctx context.Context, username string, tx db.DbConn) (*models.User, error) {
	db := db.GetDb(s.db, tx)
	row := db.QueryRow(ctx, "SELECT * FROM users WHERE username = $1", username)
	foundUser := &models.User{}
	if err := util.ScanRowIntoStruct(row, foundUser); err != nil {
		return nil, errors.Wrap(err, "userStore.GetUserByUsername")
	}

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

func (s *userStore) GenerateRecoveryCode(ctx context.Context, userId int, tx db.DbConn) (*models.RecoveryCode, error) {
	db := db.GetDb(s.db, tx)

	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "userStore.GenerateRecoveryCode")
	}

	code := base64.URLEncoding.EncodeToString(bytes)[:8]

	row := db.QueryRow(ctx, "INSERT INTO recovery_codes (user_id, code) VALUES ($1, $2) RETURNING *", userId, code)

	recoveryCode := &models.RecoveryCode{}
	if err := util.ScanRowIntoStruct(row, recoveryCode); err != nil {
		return nil, errors.Wrap(err, "userStore.GenerateRecoveryCode")
	}

	return recoveryCode, nil
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

func (s *userStore) GetRecoveryCodesByUserId(ctx context.Context, userId int, tx db.DbConn) ([]*models.RecoveryCode, error) {
	db := db.GetDb(s.db, tx)

	rows, err := db.Query(ctx, "SELECT * FROM recovery_codes WHERE user_id = $1", userId)
	if err != nil {
		return nil, errors.Wrap(err, "userStore.GetRecoveryCodesByUserId")
	}

	codes := make([]*models.RecoveryCode, 0, 10)

	for rows.Next() {
		recoveryCode := &models.RecoveryCode{}
		if err := util.ScanRowIntoStruct(rows, recoveryCode); err != nil {
			return nil, errors.Wrap(err, "userStore.GetRecoveryCodesByUserId")
		}
		codes = append(codes, recoveryCode)
	}

	return codes, nil
}

func (s *userStore) RedeemRecoveryCode(ctx context.Context, id int, tx db.DbConn) error {
	db := db.GetDb(s.db, tx)

	if _, err := db.Exec(ctx, "UPDATE recovery_codes SET is_redeemed = $1 WHERE id=$2", true, id); err != nil {
		return errors.Wrap(err, "userStore.RedeemRecoveryCode")
	}

	return nil
}

func (s *userStore) CreateLoginSession(ctx context.Context, userId int, tx db.DbConn) (int, error) {
	db := db.GetDb(s.db, tx)

	row := db.QueryRow(ctx, "INSERT INTO login_sessions (user_id, expiration_timestamp) VALUES ($1, $2) RETURNING id", userId, time.Now().Add(1*time.Minute))

	var loginSessionId int
	if err := row.Scan(&loginSessionId); err != nil {
		return -1, errors.Wrap(err, "userStore.CreateLoginSession")
	}

	return loginSessionId, nil
}

func (s *userStore) GetLoginSessionById(ctx context.Context, loginSessionId int, tx db.DbConn) (*models.LoginSession, error) {
	db := db.GetDb(s.db, tx)

	row := db.QueryRow(ctx, "SELECT * FROM login_sessions WHERE id= $1", loginSessionId)

	loginSession := &models.LoginSession{}
	if err := util.ScanRowIntoStruct(row, loginSession); err != nil {
		return nil, errors.Wrap(err, "userStore.GetLoginSessionById")
	}

	return loginSession, nil
}

func (s *userStore) DeleteLoginSessionByUserId(ctx context.Context, userId int, tx db.DbConn) error {
	db := db.GetDb(s.db, tx)

	if _, err := db.Exec(ctx, "DELETE FROM login_sessions WHERE user_id= $1", userId); err != nil {
		return errors.Wrap(err, "userStore.DeleteLoginSessionByUserId")
	}

	return nil
}

func (s *userStore) UpdatePassword(ctx context.Context, password string, userId int, tx db.DbConn) error {
	db := db.GetDb(s.db, tx)

	if _, err := db.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", password, userId); err != nil {
		return errors.Wrap(err, "userStore.UpdatePassword")
	}

	return nil
}
