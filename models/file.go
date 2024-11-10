package models

import "time"

type File struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	RepoID      int        `json:"repo_id" gorm:"not null"`
	UserID      int        `json:"user_id" gorm:"not null"`
	Filename    string     `json:"filename" gorm:"not null"`
	URL         string     `json:"url" gorm:"not null"`
	HashValue   string     `json:"hash_value"`
	RawFilename string     `json:"raw_filename"`
	Filesize    uint       `json:"filesize"`
	Width       uint       `json:"width"`
	Height      uint       `json:"height"`
	Mime        string     `json:"mime"`
	Filetype    uint8      `json:"filetype" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Repository  Repository `json:"-" gorm:"foreignKey:RepoID"`
	User        User       `json:"-" gorm:"foreignKey:UserID"`
}

type FileResponse struct {
	ID          int       `json:"id"`
	Filename    string    `json:"filename"`
	URL         string    `json:"url"`
	RawFilename string    `json:"raw_filename"`
	Filesize    uint      `json:"filesize"`
	Width       uint      `json:"width,omitempty"`
	Height      uint      `json:"height,omitempty"`
	Mime        string    `json:"mime"`
	CreatedAt   time.Time `json:"created_at"`
}
