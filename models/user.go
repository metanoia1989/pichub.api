package models

import (
	"time"
)

// UserResponse 用于API响应的用户信息结构体
type UserBaseInfo struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	Nickname  string     `json:"nickname"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	UserType  uint8      `json:"user_type"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// user 表结构
type User struct {
	UserBaseInfo `json:",inline" gorm:"embedded"` // 添加 inline 和 embedded 标签
	PasswordHash string                           `json:"-"`
}

// ToBaseInfo 将 User 转换为 UserBaseInfo
func (u *User) ToBaseInfo() *UserBaseInfo {
	return &UserBaseInfo{
		ID:        u.ID,
		Nickname:  u.Nickname,
		Username:  u.Username,
		Email:     u.Email,
		UserType:  u.UserType,
		DeletedAt: u.DeletedAt,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// 其他结构体
type RegisterRequest struct {
	Username string `json:"username" form:"username" label:"用户名" binding:"required,min=3,max=32"`
	Nickname string `json:"nickname" form:"nickname" label:"昵称" binding:"required,min=2,max=32"`
	Email    string `json:"email" form:"email" label:"邮箱" binding:"required,email"`
	Password string `json:"password" form:"password" label:"密码" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" form:"username" label:"用户名" binding:"required"`
	Password string `json:"password" form:"password" label:"密码" binding:"required"`
}
