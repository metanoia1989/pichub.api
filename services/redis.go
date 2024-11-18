package services

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"pichub.api/config"
)

// RedisService 导出的Redis服务实例
var RedisService *redis.Client

// InitRedis 初始化Redis服务
func InitRedis() error {
	RedisService = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.RedisConfig(),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
	})

	// 测试连接是否成功
	err := RedisService.Ping(context.Background()).Err()

	return err
}
