package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"pichub.api/config"
)

// RedisService 导出的Redis服务实例
var RedisService *redis.Client

const (
	// TokenKeyPrefix Redis中GitHub token的key前缀
	TokenKeyPrefix = "github:token:"
)

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

// GetCachedToken 从Redis获取缓存的token
func GetCachedToken(userID int) (string, error) {
	key := fmt.Sprintf("%s%d", TokenKeyPrefix, userID)
	return RedisService.Get(context.Background(), key).Result()
}

// SetCachedToken 设置token缓存
func SetCachedToken(userID int, token string) error {
	key := fmt.Sprintf("%s%d", TokenKeyPrefix, userID)
	// 设置token缓存，过期时间设为24小时
	return RedisService.Set(context.Background(), key, token, 1*time.Hour).Err()
}

// DeleteCachedToken 删除token缓存
func DeleteCachedToken(userID int) error {
	key := fmt.Sprintf("%s%d", TokenKeyPrefix, userID)
	return RedisService.Del(context.Background(), key).Err()
}
