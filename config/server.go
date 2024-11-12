package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Host                 string `mapstructure:"SERVER_HOST"`
	Port                 string `mapstructure:"SERVER_PORT"`
	Secret               string `mapstructure:"JWT_SECRET"`
	Debug                bool   `mapstructure:"APP_DEBUG"`
	AllowedHosts         string `mapstructure:"ALLOWED_HOSTS"`
	LimitCountPerRequest int64
}

func (s *ServerConfiguration) ServerConfig() string {
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")

	appServer := fmt.Sprintf("%s:%s", viper.GetString("SERVER_HOST"), viper.GetString("SERVER_PORT"))
	log.Print("Server Running at :", appServer)
	return appServer
}
