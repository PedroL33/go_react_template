package config

import (
	"log"
	"os"

	// "strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	Host        string
	Port        string
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	config := &Config{
		DatabaseUrl: os.Getenv("DB_CONNECTION_STRING"),
		Host:        os.Getenv("HOST"),
		Port:        os.Getenv("PORT"),
	}

	return config, nil
}

// // Helper function to get a string value from environment variables with default fallback
// func getEnv(key, defaultValue string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return defaultValue
// }

// // Helper function to get a boolean value from environment variables with default fallback
// func getEnvAsBool(key string, defaultValue bool) bool {
// 	if value, exists := os.LookupEnv(key); exists {
// 		b, err := strconv.ParseBool(value)
// 		if err == nil {
// 			return b
// 		}
// 	}
// 	return defaultValue
// }
