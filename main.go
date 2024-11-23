package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"example/dashboard/api"
	"example/dashboard/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file source driver
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	err := godotenv.Load()

	config, err := config.InitConfig()

	// Database connection
	dataSourceName := os.Getenv("DB_CONNECTION_STRING")

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatalf("Could not start SQL driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "dashboard", driver)
	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	srv := api.NewServer(logger, config, db)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}
	// Start server
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
