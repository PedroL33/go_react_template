package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl  string
	Host         string
	Port         string
	JwtSecret    string
	Env          string
	MigrationUrl string
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	config := &Config{
		DatabaseUrl:  os.Getenv("DATABASE_URL"),
		Host:         os.Getenv("HOST"),
		Port:         os.Getenv("PORT"),
		JwtSecret:    os.Getenv("JWT_SECRET"),
		Env:          os.Getenv("ENVIRONMENT"),
		MigrationUrl: os.Getenv("MIGRATION_URL"),
	}

	return config, nil
}
