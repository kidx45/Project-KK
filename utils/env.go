package utils

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_URL                 string        `mapstructure:"DB_URL"`
	PORT                   string        `mapstructure:"PORT"`
	GRPC_PORT              string        `mapstructure:"GRPC_PORT"`
	DB_DRIVER_NAME         string        `mapstructure:"DB_DRIVER_NAME"`
	SYMMETRIC_SECRET_KEY   string        `mapstructure:"SYMMETRIC_SECRET_KEY"`
	ACCESS_TOKEN_DURATION  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	REFRESH_TOKEN_DURATION time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadEnv(path string) (Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	// Load the .env file from the current directory
	err := viper.ReadInConfig()
	if err != nil {
		return Config{
			DB_URL:               os.Getenv("DB_URL"),
			PORT:                 os.Getenv("PORT"),
			DB_DRIVER_NAME:       os.Getenv("DB_DRIVER_NAME"),
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
