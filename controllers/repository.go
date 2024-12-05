package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pichub.api/constants"
	"pichub.api/infra/logger"
	"pichub.api/models"
	"pichub.api/pkg/validator"
	"pichub.api/routers/middleware"
	"pichub.api/services"
)

// AddRepository 添加新的GitHub仓库
func AddRepository(c *gin.Context) {
	logger.Infof("Received request body: %v\n", c.Request.Body)
	userID, _ := middleware.GetCurrentUser(c)

	var req models.AddRepositoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.TranslateErr(err)})
		return
	}

	logger.Infof("AddRepository%s", "测试")

	if req.RepoBranch == "" {
		req.RepoBranch = constants.DefaultRepoBranch
	}

	repository, err := services.RepositoryService.AddRepository(userID, req.RepoName, req.RepoURL, req.RepoBranch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	_, err = services.RepositoryService.InitRepository(userID, repoID)
	if err != nil {
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

	repositories, err := services.RepositoryService.ListRepositories(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repositories"})
		return
	}

	var response []models.RepositoryResponse
	for _, repo := range repositories {
		response = append(response, models.RepositoryResponse{
			ID:         repo.ID,
			RepoName:   repo.RepoName,
			RepoBranch: repo.RepoBranch,
			RepoURL:    repo.RepoURL,
			CreatedAt:  repo.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"repositories": response,
	})
}

// GetRepository 获取仓库信息
func GetRepository(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)
	repoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	repository, err := services.RepositoryService.GetRepository(userID, repoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repository": models.RepositoryResponse{
			ID:         repository.ID,
			RepoName:   repository.RepoName,
			RepoBranch: repository.RepoBranch,
			RepoURL:    repository.RepoURL,
			CreatedAt:  repository.CreatedAt,
		},
	})
}

// UpdateRepository 更新仓库信息
func UpdateRepository(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)
	repoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	var req models.UpdateRepositoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.TranslateErr(err)})
		return
	}

	if err := services.RepositoryService.UpdateRepository(userID, repoID, req.RepoURL, req.RepoBranch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update repository"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Repository updated successfully",
	})
}

// DeleteRepository 删除仓库
func DeleteRepository(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)
	repoID, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	if err := services.RepositoryService.DeleteRepository(userID, repoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete repository"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repository deleted successfully"})
}
