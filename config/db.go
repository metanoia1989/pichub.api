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

	masterDBDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		masterDBUser, masterDBPassword, masterDBHost, masterDBPort, masterDBName,
	)

	return masterDBDSN
}
