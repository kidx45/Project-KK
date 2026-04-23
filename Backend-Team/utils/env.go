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
    // Load the .env file from the current directory
    err := godotenv.Load("../../.env.local.development")
    if err != nil {
        return Config{}, err
    }

	AppConfig = Config{
		DB_URL: os.Getenv("DB_URL"),
		DB_DRIVER_NAME: os.Getenv("DB_DRIVER_NAME"),
		PORT: os.Getenv("PORT"),
	}

	return AppConfig, nil
}
