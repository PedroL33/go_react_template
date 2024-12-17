package server

import (
	"example/dashboard/api/db"
	"example/dashboard/config"
	"example/dashboard/util"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
}

func NewServer(
	logger util.Logger,
	config *config.Config,
	db db.DbConn,
) http.Handler {

	mux := mux.NewRouter()

	// //var handler http.Handler = mux

	AddRoutes(mux, logger, config, db)
	return mux
}
