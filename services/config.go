package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"pichub.api/infra/database"
	"pichub.api/infra/logger"
	"pichub.api/models"
	"pichub.api/repository"
)

type ConfigServiceImpl struct{}

var ConfigService = &ConfigServiceImpl{}

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
	err := database.DB.Where("type = ? AND name = ? AND user_id = ?", configType, name, userID).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新记录
			config = models.Config{
				UserID: userID,
				Type:   configType,
				Name:   name,
				Value:  valueStr,
			}
			return repository.Save(&config).(error)
		}
		return err
	}

	// 更新现有记录
	config.Value = valueStr
	return database.DB.Save(&config).Error
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
