package handlers

import (
	"example/dashboard/api/base"
	"example/dashboard/api/models"
	"example/dashboard/api/users"
	"example/dashboard/config"
	"example/dashboard/util"
	"net/http"
)

type usersHandlers struct {
	controller users.Controller
	cfg        *config.Config
	logger     util.Logger
}

func NewUsersHandlers(cfg *config.Config, controller users.Controller, logger util.Logger) users.Handler {
	return &usersHandlers{cfg: cfg, controller: controller, logger: logger}
}

func (h *usersHandlers) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := &models.User{}
	err := util.ReadRequest(ctx, r, user)

	if err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	createdUser, err := h.controller.CreateUser(ctx, user)

	if err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	h.logger.HttpSuccess(r)
	util.Encode(w, http.StatusCreated, createdUser)
}

func (h *usersHandlers) Login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user := &models.User{}
	var err error

	if err = util.ReadRequest(ctx, r, user); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	var loggedInUser *models.UserWithToken
	if loggedInUser, err = h.controller.Login(ctx, user); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	util.Encode(w, http.StatusAccepted, loggedInUser)
	h.logger.HttpSuccess(r)
}
