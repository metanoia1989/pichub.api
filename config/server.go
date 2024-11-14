package config

import (
	"fmt"

	"pichub.api/infra/logger"
)

type ServerConfiguration struct {
	Host                 string `mapstructure:"SERVER_HOST"`
	Port                 string `mapstructure:"SERVER_PORT"`
	Secret               string `mapstructure:"JWT_SECRET"`
	Debug                bool   `mapstructure:"APP_DEBUG"`
	AllowedHosts         string `mapstructure:"ALLOWED_HOSTS"`
	LimitCountPerRequest int    `mapstructure:"LIMIT_COUNT_PER_REQUEST"`
	FrontendUrl          string `mapstructure:"FRONTEND_URL"`
}

func (s *ServerConfiguration) ServerConfig() string {
	appServer := fmt.Sprintf("%s:%s", s.Host, s.Port)
	logger.Infof("Server Running at : %s", appServer)
	return appServer
}

func (s *ServerConfiguration) GetFrontendUrl() string {
	if s.FrontendUrl == "" {
		return fmt.Sprintf("http://%s:%s", s.Host, s.Port)
	}
	return s.FrontendUrl
}
