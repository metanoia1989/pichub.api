package models

import "time"

// config 表结构
type Config struct {
	ID        int       `gorm:"primaryKey;column:id"`
	UserID    int       `gorm:"column:user_id;default:0"`
	Type      string    `gorm:"column:type;size:24"`
	Name      string    `gorm:"column:name;size:32"`
	Value     string    `gorm:"column:value;size:500"`
	Remark    string    `gorm:"column:remark;size:50"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// func (Config) TableName() string {
// 	return "config"
// }
