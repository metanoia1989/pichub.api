package config

import (
	"github.com/spf13/viper"
	"pichub.api/infra/logger"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

// SetupConfig configuration
func SetupConfig() (*Configuration, error) {
	conf := &Configuration{}

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error to reading config file, %s", err)
		return nil, err
	}

	err := viper.Unmarshal(conf)
	if err != nil {
		logger.Errorf("error to decode, %v", err)
		return nil, err
	}

	return conf, nil
}
