package controller

import (
	"context"
	"example/dashboard/api/models"
	"example/dashboard/api/users"
	"example/dashboard/config"
	http_errors "example/dashboard/errors"
	"example/dashboard/util"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type usersController struct {
	cfg        *config.Config
	usersStore users.Store
	logger     util.Logger
}

func NewUsersController(cfg *config.Config, usersStore users.Store, logger util.Logger) users.Controller {
	return &usersController{cfg: cfg, usersStore: usersStore, logger: logger}
}

func (uc *usersController) CreateUser(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	var err error
	var u *models.User
	if u, err = uc.usersStore.GetUserByEmail(ctx, user.Email); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.CreateUser"))
		}
	}

	if u != nil {
		return nil, http_errors.NewHttpError(http.StatusBadRequest, "Email already exists.", errors.Wrap(err, "UsersController.CreateUser"))
	}

	var userWithToken *models.UserWithToken

	user.PrepareCreate()
	userWithToken, err = uc.usersStore.CreateUser(ctx, user, func() (string, error) {
		return util.CreateToken(uc.cfg, user)
	})

	if err != nil {
		return nil, http_errors.NewBadRequestError(errors.Wrap(err, "UsersController.CreateUser"))
	}

	userWithToken.User.Sanitize()

	return userWithToken, nil
}

func (uc *usersController) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	var foundUser *models.User
	var err error

	if foundUser, err = uc.usersStore.GetUserByEmail(ctx, user.Email); err != nil {
		return nil, http_errors.InvalidCredentialsError(errors.Wrap(err, "UsersController.Login"))
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, http_errors.InvalidCredentialsError(errors.Wrap(err, "UsersController.Login"))
	}

	var token string
	if token, err = util.CreateToken(uc.cfg, foundUser); err != nil {
		return nil, http_errors.NewInternalServerError(errors.Wrap(err, "UsersController.Login"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil

}
