package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pichub.api/models"
	"pichub.api/services"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UserService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 发送激活邮件
	if err := services.EmailService.SendActivationEmail(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send activation email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful. Please check your email to activate your account.",
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := services.UserService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func ActivateAccount(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Activation token is required"})
		return
	}

	if err := services.UserService.ActivateAccount(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account activated successfully"})
}
