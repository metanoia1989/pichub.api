package models

import "time"

// backup_record 表结构
type BackupRecord struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	RepoID     int       `json:"repo_id" gorm:"not null"`
	BackupPath string    `json:"backup_path" gorm:"not null"`
	BackupDate time.Time `json:"backup_date" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}

// 其他结构体

type BackupConfig struct {
	BackupDir    string
	DatabaseName string
	Username     string
	Password     string
	Host         string
	Port         string
}
