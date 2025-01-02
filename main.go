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
	conf, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	logger := logger.NewLogger()

	dbConn, err := db.NewDbConn(ctx, conf.DatabaseUrl)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := dbConn.Close(ctx); err != nil {
			logger.Error(err)
		}
	}()
	logger.Info("Database connection established.")

	db.Migrate(conf, dbConn, logger)
	logger.Info("Migrations performed and up to date.")

	db.Seed(conf, dbConn, logger)
	logger.Info("Database seed process complete.")

	srv := server.NewServer(logger, conf, dbConn)
	srv.Run()
}
