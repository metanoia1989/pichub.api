package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pichub.api/controllers"
	"pichub.api/infra/logger"
	"pichub.api/routers/middleware"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": "Route Not Found"})
	})

	// API v1 路由组
	v1 := route.Group("/api/v1")
	{
		v1.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })

		// 公开路由
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.GET("/activate", controllers.ActivateAccount)
			auth.POST("/activate", controllers.ActivateAccount)
		}

		// 需要认证的路由
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			// 这里添加需要认证的路由
			// 例如：
			// protected.GET("/user/profile", controllers.GetUserProfile)
			// protected.POST("/repositories", controllers.AddRepository)

			user := protected.Group("/user")
			{
				user.GET("/profile", controllers.GetUserProfile)
				user.POST("/profile", controllers.UpdateUserProfile)
				user.GET("/github_token", controllers.CheckGithubToken)
				user.POST("/github_token", controllers.UpdateGithubToken)
				user.POST("/email", controllers.UpdateEmail)
				user.POST("/email/verification", controllers.SendEmailVerification)
			}

			repo := protected.Group("/repositories")
			{
				logger.Infof("Registering repository routes\n")
				repo.GET("", controllers.ListRepositories)
				repo.GET("/:id", controllers.GetRepository)
				repo.POST("", controllers.AddRepository)
				repo.POST("/:id", controllers.UpdateRepository)
				repo.POST("/:id/init", controllers.InitRepository)
				repo.POST("/:id/delete", controllers.DeleteRepository)
			}

			// 在 protected 路由组中添加
			files := protected.Group("/files")
			{
				files.GET("", controllers.ListFiles)
				files.POST("/upload", controllers.UploadFile)
				files.POST("/uploadStream/:repo_id", controllers.UploadStream)
				files.POST("/delete", controllers.DeleteFile)
			}

			config := protected.Group("/config")
			{
				config.GET("/all", controllers.GetAllConfig)
				config.POST("/create", controllers.CreateConfig)
				config.POST("/update", controllers.UpdateConfig)
				config.POST("/delete/:id", controllers.DeleteConfig)
			}
		}

		// 其他的公开路由
		// 处理 github webhook 请求
		v1.POST("/webhook/github", controllers.GithubWebhook)
	}
}
