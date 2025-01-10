package server

import (
	"context"
	"example/dashboard/api/db"
	"example/dashboard/api/middleware"
	"example/dashboard/api/users/controller"
	"example/dashboard/api/users/handlers"
	"example/dashboard/api/users/store"
	"example/dashboard/config"
	"example/dashboard/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	logger logger.Logger
	conf   *config.AppConfig
	db     db.DbConn
}

func NewServer(
	logger logger.Logger,
	conf *config.AppConfig,
	db db.DbConn,
) *Server {
	return &Server{logger: logger, conf: conf, db: db}
}

func (s *Server) Run() {

	router := mux.NewRouter()
	s.MapRoutes(router)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(s.conf.Host, s.conf.Port),
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM) // Listen for SIGINT (Ctrl+C) and SIGTERM (e.g., from Kubernetes)

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		s.logger.Info("Starting server", "address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error(err)
			stop <- syscall.SIGTERM // Force shutdown if the server fails
		}
	}()

	// Wait for shutdown signal
	<-stop
	s.logger.Info("Shutting down server...")

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		s.logger.Error(err)
	} else {
		s.logger.Info("Server stopped gracefully")
	}
}

func (s *Server) MapRoutes(router *mux.Router) {
	txnManager := db.NewTransactionManager(s.db)

	usersStore := store.NewUsersStore(s.db)
	usersController := controller.NewUsersController(s.conf, usersStore, txnManager, s.logger)
	usersHandler := handlers.NewUsersHandlers(s.conf, usersController, s.logger)

	mw := middleware.NewMiddlewareManager(s.conf, s.logger, usersStore)

	// Apply cors mw
	router.Use(mw.Cors)

	// Map routes to router
	handlers.MapUsersRoutes(usersHandler, router, mw)
}
