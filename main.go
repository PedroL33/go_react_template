package main

import (
	"context"
	"example/dashboard/api/db"
	"example/dashboard/api/server"
	"example/dashboard/config"
	"example/dashboard/util"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()
	conf, err := config.InitConfig()

	logger := util.NewLogger()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

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
	srv := server.NewServer(logger, conf, dbConn)
	srv.Run()
}
