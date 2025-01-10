package db

import (
	"context"
	"example/dashboard/api/models"
	"example/dashboard/config"
	"example/dashboard/logger"
)

func Seed(conf *config.AppConfig, db DbConn, logger logger.Logger) {
	ctx := context.Background()
	row := db.QueryRow(ctx, "SELECT COUNT(*) FROM users")

	var count int
	if err := row.Scan(&count); err != nil {
		logger.Error(err)
		return
	}

	if count > 0 {
		return
	}

	user := &models.User{
		Username: "admin",
		Password: "password",
	}

	user.HashPassword()

	if _, err := db.Exec(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password); err != nil {
		logger.Error(err)
		return
	}
}
