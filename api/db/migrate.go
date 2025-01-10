package db

import (
	"example/dashboard/config"
	"example/dashboard/logger"

	"github.com/golang-migrate/migrate/v4"
)

func Migrate(conf *config.AppConfig, db DbConn, logger logger.Logger) {
	m, err := migrate.New(
		"file://./migrations",
		conf.MigrationUrl,
	)

	if err != nil {
		logger.Error(err)
		return
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error(err)
		return
	}
}
