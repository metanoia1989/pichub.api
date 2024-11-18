package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"pichub.api/infra/database"
	"pichub.api/infra/logger"
	"pichub.api/models"
	"pichub.api/pkg/utils"
	"pichub.api/repository"
)

type ConfigServiceImpl struct{}

var ConfigService = &ConfigServiceImpl{}

const (
	// TokenKeyPrefix Redis中GitHub token的key前缀
	TokenKeyPrefix = "github:token:"
	// CDNHostKeyPrefix Redis中CDN域名的key前缀
	CDNHostKeyPrefix = "file:cdn_host:"
)

// Set 设置配置
func (s *ConfigServiceImpl) Set(configType, name string, value interface{}, userID int) error {
	var valueStr string
	switch v := value.(type) {
	case string:
		valueStr = v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		valueStr = fmt.Sprintf("%v", v)
	default:
		bytes, err := json.Marshal(value)
		if err != nil {
			logger.Errorf("marshal config value error: %v", err)
			return err
		}
		valueStr = string(bytes)
	}

	var config models.Config
	if err := database.DB.Where("type = ? AND name = ? AND user_id = ?", configType, name, userID).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新记录
			config = models.Config{
				UserID: userID,
				Type:   configType,
				Name:   name,
				Value:  valueStr,
			}
			logger.Infof("Create config: %v", config)
			if err := repository.Save(&config); err != nil {
				return err
			}
		}
		return errors.New(err.Error())
	}

	// 更新现有记录
	config.Value = valueStr
	if err := database.DB.Save(&config).Error; err != nil {
		return err
	}
	return nil
}

// Get 获取配置
func (s *ConfigServiceImpl) Get(configType, name string, userID int) (interface{}, error) {
	var config models.Config
	err := database.DB.Where("type = ? AND name = ? AND user_id = ?", configType, name, userID).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 从 viper 获取默认值
			viperKey := configType + "." + name
			if viper.IsSet(viperKey) {
				return viper.Get(viperKey), nil
			}
			return nil, nil
		}
		return nil, err
	}

	// 尝试解析JSON
	var jsonValue interface{}
	if err := json.Unmarshal([]byte(config.Value), &jsonValue); err == nil {
		return jsonValue, nil
	}

	// 如果不是JSON，返回原始字符串
	return config.Value, nil
}

// GetByType 获取指定类型的所有配置
func (s *ConfigServiceImpl) GetByType(configType string, userID int) (map[string]interface{}, error) {
	var configs []models.Config
	err := database.DB.Where("type = ? AND user_id = ?", configType, userID).Find(&configs).Error
	if err != nil {
		return nil, err
	}

	configMap := make(map[string]interface{})
	for _, config := range configs {
		var value interface{}
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			value = config.Value
		}
		configMap[config.Name] = value
	}

	// 如果数据库中没有数据，从 viper 获取
	if len(configMap) == 0 && viper.IsSet(configType) {
		return viper.GetStringMap(configType), nil
	}

	return configMap, nil
}

// GetGithubToken 获取用户的GitHub token
func (s *ConfigServiceImpl) GetGithubToken(userID int) (string, error) {
	// 先尝试从Redis缓存获取
	key := fmt.Sprintf("%s%d", TokenKeyPrefix, userID)
	if token, err := RedisService.Get(context.Background(), key).Result(); err == nil {
		return token, nil
	}

	// 缓存不存在，从数据库获取
	token, err := s.Get("github", "token", userID)
	if err != nil {
		return "", err
	}

	if utils.IsEmpty(token) {
		return "", errors.New("github token not found")
	}

	// 设置缓存，1小时过期
	if err := RedisService.Set(context.Background(), key, token.(string), 1*time.Hour).Err(); err != nil {
		logger.Warnf("Failed to cache token: %v", err)
	}

	return token.(string), nil
}

// SetGithubToken 设置用户的GitHub token
func (s *ConfigServiceImpl) SetGithubToken(userID int, token string) error {
	// 先更新数据库
	if err := s.Set("github", "token", token, userID); err != nil {
		return err
	}

	// 删除缓存
	key := fmt.Sprintf("%s%d", TokenKeyPrefix, userID)
	if err := RedisService.Del(context.Background(), key).Err(); err != nil {
		logger.Warnf("Failed to delete token cache: %v", err)
	}

	return nil
}

// GetFileCDNHostname 获取用户的CDN域名
func (s *ConfigServiceImpl) GetFileCDNHostname(userID int) string {
	// 尝试从缓存获取
	key := fmt.Sprintf("%s%d", CDNHostKeyPrefix, userID)
	if hostname, err := RedisService.Get(context.Background(), key).Result(); err == nil {
		return hostname
	}

	// 从配置获取
	hostname, err := s.Get("file", "cdn_host", userID)
	if err != nil || utils.IsEmpty(hostname) {
		return ""
	}

	// 设置缓存，1小时过期
	_ = RedisService.Set(context.Background(), key, hostname.(string), 1*time.Hour)

	return hostname.(string)
}

// SetFileCDNHostname 设置用户的CDN域名
func (s *ConfigServiceImpl) SetFileCDNHostname(userID int, hostname string) error {
	// 更新配置
	if err := s.Set("file", "cdn_host", hostname, userID); err != nil {
		return err
	}

	// 删除缓存
	key := fmt.Sprintf("%s%d", CDNHostKeyPrefix, userID)
	if err := RedisService.Del(context.Background(), key).Err(); err != nil {
		logger.Warnf("Failed to delete cdn hostname cache: %v", err)
	}

	return nil
}
