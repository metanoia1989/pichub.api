package models

import (
	"strings"
	"time"
)

// repository 表结构
type Repository struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	UserID     int       `json:"user_id" gorm:"not null"`
	RepoName   string    `json:"repo_name" gorm:"not null"`
	RepoURL    string    `json:"repo_url" gorm:"not null"`
	RepoBranch string    `json:"repo_branch" gorm:"not null;default:master"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
}

func (r *Repository) GetRepositoryName() string {
	repoPath := strings.Split(r.RepoURL, "/")
	if len(repoPath) > 0 {
		return repoPath[len(repoPath)-1]
	}
	return ""
}

// 其他结构体

type AddRepositoryRequest struct {
	RepoName   string `json:"repo_name" form:"repo_name" label:"仓库名称" binding:"required"`
	RepoURL    string `json:"repo_url" form:"repo_url" label:"仓库URL" binding:"required"`
	RepoBranch string `json:"repo_branch" form:"repo_branch" label:"仓库分支"`
}

type RepositoryResponse struct {
	ID         int       `json:"id"`
	RepoName   string    `json:"repo_name"`
	RepoURL    string    `json:"repo_url"`
	RepoBranch string    `json:"repo_branch"`
	CreatedAt  time.Time `json:"created_at"`
}

type UpdateRepositoryRequest struct {
	RepoName   string `json:"repo_name" form:"repo_name" label:"仓库名称" binding:"required"`
	RepoURL    string `json:"repo_url" form:"repo_url" label:"仓库URL" binding:"required"`
	RepoBranch string `json:"repo_branch" form:"repo_branch" label:"仓库分支" binding:"required"`
}
