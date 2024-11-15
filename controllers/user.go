package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pichub.api/models"
	"pichub.api/pkg/validator"
	"pichub.api/routers/middleware"
	"pichub.api/services"
)

// GetUserProfile 获取用户信息
func GetUserProfile(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)
	user, err := services.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.ToBaseInfo()})
}

// UpdateUserProfile 更新用户信息
func UpdateUserProfile(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var req models.UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.TranslateErr(err)})
		return
	}

	if err := services.UserService.UpdateProfile(userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的用户信息
	user, err := services.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user.ToBaseInfo(),
	})
}

// CheckGithubToken 检查github token
func CheckGithubToken(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	hasGithubToken, err := services.UserService.HasGithubToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check github token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_github_token": hasGithubToken})

}

// UpdateGithubToken 更新 github token
func UpdateGithubToken(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var req models.UpdateGithubTokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validator.TranslateErr(err),
		})
		return
	}

	// 更新 token
	if err := services.UserService.UpdateGithubToken(userID, req.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Github token updated successfully",
	})
}

// SendEmailVerification 发送邮箱验证码
func SendEmailVerification(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var req models.SendEmailVerificationRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.TranslateErr(err)})
		return
	}

	if err := services.UserService.SendEmailVerificationCode(userID, req.NewEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification code sent successfully",
	})
}

// UpdateEmail 确认验证码并更新邮箱
func UpdateEmail(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var req models.UpdateEmailRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.TranslateErr(err)})
		return
	}

	if err := services.UserService.UpdateEmail(userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的用户信息
	user, err := services.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email updated successfully",
		"user":    user.ToBaseInfo(),
	})
}
