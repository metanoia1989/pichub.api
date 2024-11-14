package models

import "time"

// repository 表结构
type Repository struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"not null"`
	RepoName  string    `json:"repo_name" gorm:"not null"`
	RepoURL   string    `json:"repo_url" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
}

// 其他结构体

type AddRepositoryRequest struct {
	RepoName string `json:"repo_name" form:"repo_name" label:"仓库名称" binding:"required"`
	RepoURL  string `json:"repo_url" form:"repo_url" label:"仓库URL" binding:"required,url"`
}

type RepositoryResponse struct {
	ID        int       `json:"id"`
	RepoName  string    `json:"repo_name"`
	RepoURL   string    `json:"repo_url"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateRepositoryRequest struct {
	RepoURL string `json:"repo_url" form:"repo_url" label:"仓库URL" binding:"required,url"`
}
