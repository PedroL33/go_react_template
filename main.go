package main

import (
	"context"
	"example/dashboard/api/db"
	"example/dashboard/api/server"
	"example/dashboard/config"
	"example/dashboard/util"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	conf, err := config.InitConfig()

	// logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger := util.NewLogger()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db, err := db.NewDbConn(context.Background(), conf.DatabaseUrl)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	m, err := migrate.New(
		"file://./migrations",
		"pgx5://peter:Mynewpw1!@localhost:5432/dashboard",
	)

	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	srv := server.NewServer(logger, conf, db)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(conf.Host, conf.Port),
		Handler: srv,
	}
	// Start server
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
