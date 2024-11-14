package config

import (
	"github.com/spf13/viper"
	"pichub.api/infra/logger"
)

type Configuration struct {
	Server   ServerConfiguration   `mapstructure:",squash"`
	Database DatabaseConfiguration `mapstructure:",squash"`
	Redis    RedisConfiguration    `mapstructure:",squash"`
	Email    EmailConfiguration    `mapstructure:",squash"`
}

var Config = &Configuration{}

// SetupConfig configuration
func SetupConfig() (*Configuration, error) {

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error to reading config file, %s", err)
		return nil, err
	}

	// 打印所有环境变量，用于调试
	// logger.Infof("All settings from viper: %v", viper.AllSettings())

	// 设置默认值
	setDefaultConfig()

	err := viper.Unmarshal(Config)
	if err != nil {
		logger.Errorf("error to decode, %v", err)
		return nil, err
	}

	// 打印解析后的配置
	// logger.Infof("Loaded configuration: %+v", Config)

	// 打印邮件配置
	// logger.Infof("Email configuration: %+v", Config.Email)

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

	// 添加数据库默认值
	viper.SetDefault("DB_DRIVER", "mysql")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_DATABASE", "pichub")
	viper.SetDefault("DB_USERNAME", "root")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_PREFIX", "")
}
