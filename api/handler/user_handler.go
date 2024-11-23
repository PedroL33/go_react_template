package handler

import (
	"example/dashboard/api/controller"
	"example/dashboard/util"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	UserController controller.UserController
	Logger         *slog.Logger
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	logger := h.Logger

	var decoded userRequest
	var err error

	if decoded, err = util.Decode[userRequest](r); err != nil {
		logger.Error(err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err = h.UserController.UserStore.CreateUser(decoded.Email, decoded.Password, decoded.FirstName, decoded.LastName); err != nil {
		logger.Error(err.Error())
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
