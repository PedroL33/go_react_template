package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DatabaseUrl  string
	Host         string
	Port         string
	JwtSecret    string
	Env          string
	MigrationUrl string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

func (ac *AppConfig) Load() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	ac.DatabaseUrl = os.Getenv("DATABASE_URL")
	ac.Host = os.Getenv("HOST")
	ac.Port = os.Getenv("PORT")
	ac.JwtSecret = os.Getenv("JWT_SECRET")
	ac.Env = os.Getenv("ENVIRONMENT")
	ac.MigrationUrl = os.Getenv("MIGRATION_URL")

	return nil
}
