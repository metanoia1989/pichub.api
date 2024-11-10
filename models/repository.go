package models

import "time"

type Repository struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"not null"`
	RepoName  string    `json:"repo_name" gorm:"not null"`
	RepoURL   string    `json:"repo_url" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
}

type AddRepositoryRequest struct {
	RepoName string `json:"repo_name" binding:"required"`
	RepoURL  string `json:"repo_url" binding:"required,url"`
}

type RepositoryResponse struct {
	ID        int       `json:"id"`
	RepoName  string    `json:"repo_name"`
	RepoURL   string    `json:"repo_url"`
	CreatedAt time.Time `json:"created_at"`
}
