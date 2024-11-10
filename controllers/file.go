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

// UploadFile 处理文件上传请求
func UploadFile(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	// 获取仓库ID
	repoID, err := strconv.Atoi(c.PostForm("repo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	// 验证仓库权限
	var repo models.Repository
	if err := database.DB.Where("id = ? AND user_id = ?", repoID, userID).First(&repo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 检查是否强制上传
	isForce := c.PostForm("is_force") == "true"

	// 处理文件上传
	uploadedFile, err := services.FileService.UploadFile(file, userID, repoID, isForce)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file": models.FileResponse{
			ID:          uploadedFile.ID,
			Filename:    uploadedFile.Filename,
			URL:         uploadedFile.URL,
			RawFilename: uploadedFile.RawFilename,
			Filesize:    uploadedFile.Filesize,
			Width:       uploadedFile.Width,
			Height:      uploadedFile.Height,
			Mime:        uploadedFile.Mime,
			CreatedAt:   uploadedFile.CreatedAt,
		},
	})
}

// ListFiles 获取用户的所有文件
func ListFiles(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var files []models.File
	query := database.DB.Where("user_id = ?", userID)

	// 支持按仓库ID筛选
	if repoID := c.Query("repo_id"); repoID != "" {
		query = query.Where("repo_id = ?", repoID)
	}

	if err := query.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}

	var response []models.FileResponse
	for _, file := range files {
		response = append(response, models.FileResponse{
			ID:          file.ID,
			Filename:    file.Filename,
			URL:         file.URL,
			RawFilename: file.RawFilename,
			Filesize:    file.Filesize,
			Width:       file.Width,
			Height:      file.Height,
			Mime:        file.Mime,
			CreatedAt:   file.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": response,
	})
}
