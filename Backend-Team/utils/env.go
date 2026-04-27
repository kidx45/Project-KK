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

func LoadEnv() (Config, error) {
    // Load the .env file from the current directory
    _ = godotenv.Load("../../.env.local.development")

	AppConfig := Config{
		DB_URL: os.Getenv("DB_URL"),
		DB_DRIVER_NAME: os.Getenv("DB_DRIVER_NAME"),
		PORT: os.Getenv("PORT"),
	}

	return AppConfig, nil
}
