package api

import (
	"database/sql"
	"example/dashboard/config"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
}

func NewServer(
	logger *slog.Logger,
	config *config.Config,
	db *sql.DB,
) http.Handler {

	mux := mux.NewRouter()

	// var handler http.Handler = mux

	AddRoutes(mux, logger, config, db)
	return mux
}
