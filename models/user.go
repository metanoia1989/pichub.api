package models

import (
	"time"
)

type User struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	Nickname     string     `json:"nickname"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	UserType     uint8      `json:"user_type"`
	DeletedAt    *time.Time `json:"deleted_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Nickname string `json:"nickname" binding:"required,min=2,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
