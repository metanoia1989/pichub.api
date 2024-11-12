package config

import (
	"github.com/spf13/viper"
)

type RedisConfiguration struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

func (r *RedisConfiguration) RedisConfig() string {
	return viper.GetString("REDIS_HOST") + ":" + viper.GetString("REDIS_PORT")
}
