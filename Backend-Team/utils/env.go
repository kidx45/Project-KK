package utils

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
	DB_URL string
	PORT string
	DB_DRIVER_NAME string
}

var AppConfig Config

func LoadEnv() (Config, error) {
    // Attempt to load the .env file, but ignore the error if it doesn't exist
    // This allows the app to fall back to system environment variables in CI/CD.
    _ = godotenv.Load("../../.env.local.development")

	AppConfig = Config{
		DB_URL: os.Getenv("DB_URL"),
		DB_DRIVER_NAME: os.Getenv("DB_DRIVER_NAME"),
		PORT: os.Getenv("PORT"),
	}

	return AppConfig, nil
}
