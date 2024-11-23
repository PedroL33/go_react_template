package api

import (
	"database/sql"
	"example/dashboard/api/controller"
	"example/dashboard/api/handler"
	"example/dashboard/api/store"
	"log/slog"

	"example/dashboard/config"

	"github.com/gorilla/mux"
)

func AddRoutes(
	mux *mux.Router,
	logger *slog.Logger,
	config *config.Config,
	db *sql.DB,
) {

	userStore := &store.SQLUserStore{DB: db}
	userController := controller.UserController{UserStore: userStore}
	userHandler := handler.UserHandler{UserController: userController, Logger: logger}

	mux.HandleFunc("/user", userHandler.CreateUserHandler).Methods("POST")
	mux.Handle("/", mux.NotFoundHandler)
}
