package config

type EmailConfiguration struct {
	Host        string `mapstructure:"MAIL_HOST"`
	Port        int    `mapstructure:"MAIL_PORT"`
	Username    string `mapstructure:"MAIL_USERNAME"`
	Password    string `mapstructure:"MAIL_PASSWORD"`
	FromAddress string `mapstructure:"MAIL_FROM_ADDRESS"`
	FromName    string `mapstructure:"MAIL_FROM_NAME"`
}
