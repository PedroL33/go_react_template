package db

import (
	"example/dashboard/config"
	"example/dashboard/util"

	"github.com/golang-migrate/migrate/v4"
)

func Migrate(conf *config.Config, db DbConn, logger util.Logger) {
	m, err := migrate.New(
		"file://./migrations",
		conf.MigrationUrl,
	)

	if err != nil {
		logger.Error(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error(err)
	}
}
