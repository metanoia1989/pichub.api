package config

type RedisConfiguration struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

func (r *RedisConfiguration) RedisConfig() string {
	return r.Host + ":" + r.Port
}
