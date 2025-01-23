package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"pichub.api/infra/logger"
	"pichub.api/routers/middleware"
)

/**
 * 初始化 gin 实例，然后注册路由
 */
func SetupRoute() *gin.Engine {
	environment := viper.GetBool("APP_DEBUG")
	if environment {
		logger.Infof("Setting gin mode to DebugMode")
		gin.SetMode(gin.DebugMode)
	} else {
		logger.Infof("Setting gin mode to ReleaseMode")
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	// allowedHosts := viper.GetString("ALLOWED_HOSTS")
	// router.SetTrustedProxies([]string{allowedHosts})
	router.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})
	
	// 添加详细的请求调试中间件
	router.Use(func(c *gin.Context) {
		logger.Infof("==== Request Debug Info ====")
		logger.Infof("Full URL: %v", c.Request.URL.String())
		// logger.Infof("Method: %v", c.Request.Method)
		// logger.Infof("Path: %v", c.Request.URL.Path)
		// logger.Infof("RawQuery: %v", c.Request.URL.RawQuery)
		// logger.Infof("Headers: %v", c.Request.Header)
		// logger.Infof("RemoteAddr: %v", c.Request.RemoteAddr)
		// logger.Infof("Host: %v", c.Request.Host)
		
		// // 检查是否经过代理
		// logger.Infof("X-Forwarded-For: %v", c.Request.Header.Get("X-Forwarded-For"))
		// logger.Infof("X-Real-IP: %v", c.Request.Header.Get("X-Real-IP"))
		// logger.Infof("X-Forwarded-Proto: %v", c.Request.Header.Get("X-Forwarded-Proto"))
		
		c.Next()
	})

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	RegisterRoutes(router) //routes register

	return router
}
