package store

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"example/template/rest_server/api/db"
	"example/template/rest_server/api/models"
	"example/template/rest_server/api/users"
	"example/template/rest_server/util"
	"time"

	"github.com/jackc/pgx/v5"
)

// userStore runs queries against whichever Querier it is bound to.
// Immutable — WithQuerier returns a new view rather than mutating.
type userStore struct {
	q db.Querier
}

func NewUsersStore(q db.Querier) users.Store {
	return &userStore{q: q}
}

func (s *userStore) WithQuerier(q db.Querier) users.Store {
	return &userStore{q: q}
}

func (s *userStore) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	rows, err := s.q.Query(
		ctx,
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *",
		user.Username,
		user.Password,
	)

	if err != nil {
		return nil, util.Wrap(err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[models.User])
}

func (s *userStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	rows, err := s.q.Query(ctx, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, util.Wrap(err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[models.User])
}

func (s *userStore) GetUserById(ctx context.Context, userId int) (*models.User, error) {
	rows, err := s.q.Query(ctx, "SELECT * FROM users WHERE id = $1", userId)

	if err != nil {
		return nil, util.Wrap(err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[models.User])
}

func (s *userStore) Create2faSetupSession(ctx context.Context, session *models.TwoFactorSetupSession) (*models.TwoFactorSetupSession, error) {
	rows, err := s.q.Query(
		ctx,
		"INSERT INTO two_factor_setup_sessions (user_id, secret_string, expiration_timestamp) VALUES ($1, $2, $3) RETURNING *",
		session.UserId,
		session.SecretString,
		session.Expiration,
	)

	if err != nil {
		return nil, util.Wrap(err)
	}

	session, err = pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[models.TwoFactorSetupSession])

	if err != nil {
		return nil, util.Wrap(err)
	}

	return session, nil
}

func (s *userStore) Get2faSetupSessionByUserId(ctx context.Context, userId int) (*models.TwoFactorSetupSession, error) {
	rows, err := s.q.Query(ctx, "SELECT * FROM two_factor_setup_sessions WHERE user_id = $1", userId)

	if err != nil {
		return nil, util.Wrap(err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[models.TwoFactorSetupSession])
}

func (s *userStore) EnableTwoFactorAuth(ctx context.Context, setupSession *models.TwoFactorSetupSession) error {
	if _, err := s.q.Exec(
		ctx,
		"UPDATE users SET two_factor_secret = $1, is_two_factor_enabled = $2 WHERE id=$3",
		setupSession.SecretString,
		true,
		setupSession.UserId,
	); err != nil {
		return util.Wrap(err)
	}

	return nil
}

func (s *userStore) DisableTwoFactorAuth(ctx context.Context, userId int) error {
	if _, err := s.q.Exec(
		ctx,
		"UPDATE users SET two_factor_secret = $1, is_two_factor_enabled = $2 WHERE id=$3",
		nil,
		false,
		userId,
	); err != nil {
		return util.Wrap(err)
	}

	return nil
}

func (s *userStore) Delete2faSetupSession(ctx context.Context, userId int) error {
	if _, err := s.q.Exec(
		ctx,
		"DELETE FROM two_factor_setup_sessions WHERE user_id = $1",
		userId,
	); err != nil {
		return util.Wrap(err)
	}

	return nil
}

func (s *userStore) GenerateRecoveryCode(ctx context.Context, userId int) (*models.RecoveryCode, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return nil, util.Wrap(err)
	}

	code := base64.URLEncoding.EncodeToString(bytes)[:8]

	rows, err := s.q.Query(ctx, "INSERT INTO recovery_codes (user_id, code) VALUES ($1, $2) RETURNING *", userId, code)

	if err != nil {
		return nil, util.Wrap(err)
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[models.RecoveryCode])
}

func (s *userStore) DeleteRecoveryCodes(ctx context.Context, userId int) error {
	if _, err := s.q.Exec(
		ctx,
		"DELETE FROM recovery_codes WHERE user_id = $1",
		userId,
	); err != nil {
		return util.Wrap(err)
	}

	return nil
}

func (s *userStore) GetRecoveryCodesByUserId(ctx context.Context, userId int) ([]*models.RecoveryCode, error) {
	rows, err := s.q.Query(ctx, "SELECT * FROM recovery_codes WHERE user_id = $1", userId)
	if err != nil {
		return nil, util.Wrap(err)
	}

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[models.RecoveryCode])
}

func (s *userStore) RedeemRecoveryCode(ctx context.Context, id int) error {
	if _, err := s.q.Exec(ctx, "UPDATE recovery_codes SET is_redeemed = $1 WHERE id=$2", true, id); err != nil {
		return util.Wrap(err)
	}

	return nil
}

func (s *userStore) CreateLoginSession(ctx context.Context, userId int) (int, error) {
	row := s.q.QueryRow(ctx, "INSERT INTO login_sessions (user_id, expiration_timestamp) VALUES ($1, $2) RETURNING id", userId, time.Now().Add(1*time.Minute))

	var loginSessionId int
	if err := row.Scan(&loginSessionId); err != nil {
		return -1, util.Wrap(err)
	}

	return loginSessionId, nil
}

func (s *userStore) GetLoginSessionById(ctx context.Context, loginSessionId int) (*models.LoginSession, error) {
	rows, err := s.q.Query(ctx, "SELECT * FROM login_sessions WHERE id= $1", loginSessionId)
	if err != nil {
		return nil, util.Wrap(err)
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[models.LoginSession])
}

func (s *userStore) DeleteLoginSessionByUserId(ctx context.Context, userId int) error {
	if _, err := s.q.Exec(ctx, "DELETE FROM login_sessions WHERE user_id= $1", userId); err != nil {
		return util.Wrap(err)
	}

	return nil
}

func (s *userStore) UpdatePassword(ctx context.Context, password string, userId int) error {
	if _, err := s.q.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", password, userId); err != nil {
		return util.Wrap(err)
	}

	return nil
}
