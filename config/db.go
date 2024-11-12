package config

import (
	"fmt"

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
	LogMode  bool
}

func (d *DatabaseConfiguration) DSN() string {
	masterDBName := viper.GetString("DB_DATABASE")
	masterDBUser := viper.GetString("DB_USERNAME")
	masterDBPassword := viper.GetString("DB_PASSWORD")
	masterDBHost := viper.GetString("DB_HOST")
	masterDBPort := viper.GetString("DB_PORT")

	masterDBDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		masterDBUser, masterDBPassword, masterDBHost, masterDBPort, masterDBName,
	)

	return masterDBDSN
}
