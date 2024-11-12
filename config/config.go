package config

import (
	"github.com/spf13/viper"
	"pichub.api/infra/logger"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Redis    RedisConfiguration
	Email    EmailConfiguration
}

var Config = &Configuration{}

// SetupConfig configuration
func SetupConfig() (*Configuration, error) {

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error to reading config file, %s", err)
		return nil, err
	}

	// 设置默认值
	setDefaultConfig()

	err := viper.Unmarshal(Config)
	if err != nil {
		logger.Errorf("error to decode, %v", err)
		return nil, err
	}

	return Config, nil
}

func setDefaultConfig() {
	// Server defaults
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Shanghai")

	// Redis defaults
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)
}
