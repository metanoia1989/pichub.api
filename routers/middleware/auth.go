package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"pichub.api/infra/logger"
	"pichub.api/pkg/jwt"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Infof("=== Auth middleware triggered for path: %s ===\n", c.Request.URL.Path)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// 解析token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.UserType == 1)

		c.Next()
	}
}

// GetCurrentUser 从上下文中获取当前用户ID
func GetCurrentUser(c *gin.Context) (int, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(int), true
}

func IsAdmin(c *gin.Context) bool {
	return c.GetBool("is_admin")
}
