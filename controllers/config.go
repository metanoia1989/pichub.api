package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pichub.api/infra/database"
	"pichub.api/models"
	"pichub.api/routers/middleware"
)

// GetAllConfig 获取所有配置值
func GetAllConfig(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var configs []models.Config
	if err := database.DB.Where("user_id = ? OR user_id = 0", userID).Order("user_id ASC").Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch configurations"})
		return
	}

	// 直接返回所有字段
	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// UpdateConfig 更新配置项
func UpdateConfig(c *gin.Context) {
	currentUserID, _ := middleware.GetCurrentUser(c)

	// 获取配置ID

	var req models.ConfigSetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 验证权限
	if req.UserID == 0 {
		// 系统配置需要管理员权限
		if !middleware.IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin permission required for system config"})
			return
		}
	} else if req.UserID != currentUserID {
		// 只能修改自己的配置
		c.JSON(http.StatusForbidden, gin.H{"error": "Can only modify your own config"})
		return
	}

	// 查找并更新配置
	var config models.Config
	if err := database.DB.First(&config, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}

	// 更新配置字段
	config.Type = req.Type
	config.Name = req.Name
	config.Value = req.Value
	config.UserID = req.UserID
	config.Remark = req.Remark

	if err := database.DB.Save(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration updated successfully",
		"config":  config,
	})
}

// CreateConfig 添加配置
func CreateConfig(c *gin.Context) {
	currentUserID, _ := middleware.GetCurrentUser(c)

	var req models.ConfigCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 验证权限
	if req.UserID == 0 {
		// 系统配置需要管理员权限
		if !middleware.IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin permission required for system config"})
			return
		}
	} else if req.UserID != currentUserID {
		// 只能添加自己的配置
		c.JSON(http.StatusForbidden, gin.H{"error": "Can only add config for yourself"})
		return
	}

	// 检查配置是否已存在
	var existingConfig models.Config
	if err := database.DB.Where("user_id = ? AND type = ? AND name = ?", req.UserID, req.Type, req.Name).First(&existingConfig).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Config already exists"})
		return
	}

	// 创建新配置
	config := models.Config{
		UserID: req.UserID,
		Type:   req.Type,
		Name:   req.Name,
		Value:  req.Value,
		Remark: req.Remark,
	}

	if err := database.DB.Create(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration created successfully",
		"config":  config,
	})
}

// DeleteConfig 删除配置
func DeleteConfig(c *gin.Context) {
	currentUserID, _ := middleware.GetCurrentUser(c)

	configID := c.Param("id")

	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}

	// 查找配置
	var config models.Config
	if err := database.DB.First(&config, configID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}

	// 不允许删除系统配置
	if config.UserID == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "System config cannot be deleted"})
		return
	}

	// 只能删除自己的配置
	if config.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Can only delete your own config"})
		return
	}

	// 删除配置
	if err := database.DB.Delete(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration deleted successfully",
	})
}
