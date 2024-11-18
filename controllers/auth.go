package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pichub.api/infra/logger"
	"pichub.api/models"
	"pichub.api/pkg/validator"
	"pichub.api/services"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.TranslateErr(err)})
		return
	}

	user, err := services.UserService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ID := 3
	// user, err := services.UserService.GetUserByID(ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user: " + err.Error()})
	// 	return
	// }

	// 发送激活邮件
	if err := services.EmailService.SendActivationEmail(user); err != nil {
		logger.Errorf("Failed to send activation email: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send activation email: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful. Please check your email to activate your account.",
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.TranslateErr(err)})
		return
	}

	token, user, err := services.UserService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user.ToBaseInfo(),
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
