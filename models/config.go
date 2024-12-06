package models

import (
	"time"

	"pichub.api/config"
)

// config 表结构
type Config struct {
	ID        int       `json:"id" gorm:"primaryKey;column:id"`
	UserID    int       `json:"user_id" gorm:"column:user_id;default:0"`
	Type      string    `json:"type" gorm:"column:type;size:24"`
	Name      string    `json:"name" gorm:"column:name;size:32"`
	Value     string    `json:"value" gorm:"column:value;size:500"`
	Remark    string    `json:"remark" gorm:"column:remark;size:50"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Config) TableName() string {
	return config.Config.Database.Prefix + "config"
}

type ConfigSetRequest struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	UserID int    `json:"user_id"`
	Remark string `json:"remark"`
}

type ConfigCreateRequest struct {
	UserID int    `json:"user_id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}
