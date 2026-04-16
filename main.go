package main

import (
	"context"
	"example/dashboard/api/db"
	"example/dashboard/api/server"
	"example/dashboard/config"
	"example/dashboard/logger"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()
	conf := config.NewAppConfig()
	conf.Load()

	logger := logger.NewLogger()

	pool, err := db.NewPool(ctx, conf.DatabaseUrl)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	logger.Info("Database connection pool established.")

	db.Migrate(conf, logger)
	logger.Info("Migrations performed and up to date.")

	db.Seed(conf, pool, logger)
	logger.Info("Database seed process complete.")

	srv := server.NewServer(logger, conf, pool)
	srv.Run()
}
