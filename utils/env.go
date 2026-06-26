package utils

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_URL string
	PORT string
	DB_DRIVER_NAME string
	SYMMETRIC_SECRET_KEY string
	ACCESS_TOKEN_DURATION time.Duration
}

func LoadEnv(path string) (Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
    // Load the .env file from the current directory
	err := viper.ReadInConfig()
	if err != nil {
		return Config{
			DB_URL: os.Getenv("DB_URL"),
			PORT: os.Getenv("PORT"),
			DB_DRIVER_NAME: os.Getenv("DB_DRIVER_NAME"),
			SYMMETRIC_SECRET_KEY: os.Getenv("SYMMETRIC_SECRET_KEY"),
			ACCESS_TOKEN_DURATION: func() time.Duration {
				d, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
				return d
			}(),
		}, nil
	}

	AppConfig := Config{}
	viper.Unmarshal(&AppConfig)
	return AppConfig, nil
}
