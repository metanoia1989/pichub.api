package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pichub.api/infra/database"
	"pichub.api/models"
	"pichub.api/routers/middleware"
	"pichub.api/services"
)

// AddRepository 添加新的GitHub仓库
func AddRepository(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var req models.AddRepositoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证仓库
	owner, repo, err := services.GithubService.ValidateRepository(req.RepoURL, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建仓库记录
	repository := &models.Repository{
		UserID:   userID,
		RepoName: owner + "/" + repo,
		RepoURL:  req.RepoURL,
	}

	if err := database.DB.Create(repository).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save repository"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Repository added successfully",
		"repository": models.RepositoryResponse{
			ID:        repository.ID,
			RepoName:  repository.RepoName,
			RepoURL:   repository.RepoURL,
			CreatedAt: repository.CreatedAt,
		},
	})
}

// InitRepository 初始化仓库数据
func InitRepository(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)
	repoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	// 获取仓库信息
	var repository models.Repository
	if err := database.DB.Where("id = ? AND user_id = ?", repoID, userID).First(&repository).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	// 初始化仓库数据
	if err := services.GithubService.InitializeRepository(&repository, ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Repository initialized successfully",
	})
}

// ListRepositories 获取用户的所有仓库
func ListRepositories(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var repositories []models.Repository
	if err := database.DB.Where("user_id = ?", userID).Find(&repositories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repositories"})
		return
	}

	var response []models.RepositoryResponse
	for _, repo := range repositories {
		response = append(response, models.RepositoryResponse{
			ID:        repo.ID,
			RepoName:  repo.RepoName,
			RepoURL:   repo.RepoURL,
			CreatedAt: repo.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"repositories": response,
	})
}
