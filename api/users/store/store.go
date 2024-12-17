package store

import (
	"context"
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

func (s *userStore) CreateUser(ctx context.Context, user *models.User, createToken func() (string, error)) (*models.UserWithToken, error) {

	var tx db.DbConn
	var err error
	if tx, err = s.db.Begin(ctx); err != nil {
		return nil, errors.Wrap(err, "userStore.CreateUser")
	}

	defer tx.Rollback(ctx)
	row := tx.QueryRow(
		ctx,
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *",
		user.Email,
		user.Password,
	)

	createdUser := &models.User{}
	if err := util.ScanRowIntoStruct(row, createdUser); err != nil {
		return nil, errors.Wrap(err, "userStore.CreateUser")
	}

	var token string
	if token, err = createToken(); err != nil {
		return nil, errors.Wrap(err, "userStore.CreateUser")
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "userStore.CreateUser")
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

func (s *userStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	foundUser := &models.User{}
	row := s.db.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email)

	if err := util.ScanRowIntoStruct(row, foundUser); err != nil {
		return nil, errors.Wrap(err, "userStore.GetUserByEmail")
	}

	return foundUser, nil
}
