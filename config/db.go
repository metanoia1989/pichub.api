package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
}

func DbConfiguration() string {
	masterDBName := viper.GetString("DB_DATABASE")
	masterDBUser := viper.GetString("DB_USERNAME")
	masterDBPassword := viper.GetString("DB_PASSWORD")
	masterDBHost := viper.GetString("DB_HOST")
	masterDBPort := viper.GetString("DB_PORT")
	masterDBSslMode := viper.GetString("DB_SSL_MODE")

	masterDBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		masterDBHost, masterDBUser, masterDBPassword, masterDBName, masterDBPort, masterDBSslMode,
	)

	return masterDBDSN
}
