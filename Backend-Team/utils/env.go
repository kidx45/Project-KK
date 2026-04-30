package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB_URL string
	PORT string
	DB_DRIVER_NAME string
}

func LoadEnv(path string) (Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
    // Load the .env file from the current directory
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	AppConfig := Config{}
	viper.Unmarshal(&AppConfig)
	return AppConfig, nil
}
