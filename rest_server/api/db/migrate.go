package db

import (
	"example/template/rest_server/config"
	"example/template/rest_server/logger"

	"github.com/golang-migrate/migrate/v4"
)

func Migrate(conf *config.AppConfig, logger logger.Logger) {
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
