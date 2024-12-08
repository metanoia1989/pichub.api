package config

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

type DatabaseConfiguration struct {
	Driver   string `mapstructure:"DB_DRIVER"`
	Dbname   string `mapstructure:"DB_DATABASE"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Prefix   string `mapstructure:"DB_PREFIX"`
	LogMode  bool   `mapstructure:"DB_LOG_MODE"`
}

func (d *DatabaseConfiguration) DSN() string {
	// 从 viper 获取配置的时区，如果未设置则默认使用 Asia/Shanghai
	timezone := viper.GetString("SERVER_TIMEZONE")
	if timezone == "" {
		timezone = "Asia/Shanghai"
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		d.Username, d.Password, d.Host, d.Port, d.Dbname,
		url.QueryEscape(timezone),
	)
}
