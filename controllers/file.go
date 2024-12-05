package controllers

import (
	"math"
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
	isForceParam := c.PostForm("is_force")
	isForce := isForceParam == "true" || isForceParam == "1"

	// 处理文件上传
	uploadedFile, err := services.FileService.UploadFile(file, userID, repoID, isForce)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cdnHost := services.ConfigService.GetFileCDNHostname(0)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    uploadedFile.ToResponse(cdnHost),
	})
}

// ListFiles 获取用户的所有文件
func ListFiles(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 确保页码和每页数量合理
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取仓库ID（如果有）
	repoID := 0
	if repoIDStr := c.Query("repo_id"); repoIDStr != "" {
		repoID, _ = strconv.Atoi(repoIDStr)
	}

	// 使用服务获取文件列表
	files, total, err := services.FileService.ListFiles(userID, repoID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}

	// 构建响应
	var response []models.FileResponse
	cdnHost := services.ConfigService.GetFileCDNHostname(0)
	for _, file := range files {
		response = append(response, file.ToResponse(cdnHost))
	}

	hasMore := page*pageSize < int(total)

	c.JSON(http.StatusOK, gin.H{
		"files": response,
		"pagination": gin.H{
			"current_page": page,
			"page_size":    pageSize,
			"total":        total,
			"total_pages":  int(math.Ceil(float64(total) / float64(pageSize))),
			"has_more":     hasMore,
		},
	})
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	// 定义请求体结构
	var req struct {
		ID string `json:"id"`
	}

	// c.PostForm() 只能获取 application/x-www-form-urlencoded 或 multipart/form-data 格式的数据
	// 对于 JSON 数据，需要使用 c.ShouldBindJSON() 来解析

	// 解析 JSON 请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 解析文件ID
	fileID, err := strconv.Atoi(req.ID)
	if err != nil || fileID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// 验证文件所有权并删除
	if err := services.FileService.DeleteFile(fileID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// UploadStream 处理流式文件上传请求
func UploadStream(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	// 获取仓库ID
	repoID, err := strconv.Atoi(c.Param("repo_id"))
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

	// 获取文件信息从请求头
	filename := c.GetHeader("X-File-Name")
	contentType := c.GetHeader("Content-Type")
	contentLength := c.GetHeader("Content-Length")
	isForce := c.GetHeader("X-Is-Force") == "true"

	fileSize, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Content-Length"})
		return
	}

	// 处理文件上传
	uploadedFile, err := services.FileService.UploadStream(c.Request.Body, filename, contentType, fileSize, userID, repoID, isForce)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cdnHost := services.ConfigService.GetFileCDNHostname(0)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    uploadedFile.ToResponse(cdnHost),
	})
}
