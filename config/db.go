package config

import (
	"fmt"
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
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Dbname,
	)
}
