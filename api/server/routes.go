package server

import (
	"example/dashboard/api/db"
	"example/dashboard/api/users"
	"example/dashboard/api/users/controller"
	"example/dashboard/api/users/handlers"
	"example/dashboard/api/users/store"
	"example/dashboard/config"
	"example/dashboard/util"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func AddRoutes(
	mux *mux.Router,
	logger util.Logger,
	config *config.Config,
	db db.DbConn,
) {

	usersStore := store.NewUsersStore(db)
	usersController := controller.NewUsersController(config, usersStore, logger)
	usersHandler := handlers.NewUsersHandlers(config, usersController, logger)

	users.MapRoutes(usersHandler, mux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	mux.Handle("/", mux.NotFoundHandler)
}
