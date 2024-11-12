package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"pichub.api/routers/middleware"
)

/**
 * 初始化 gin 实例，然后注册路由
 */
func SetupRoute() *gin.Engine {

	environment := viper.GetBool("APP_DEBUG")
	if environment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	allowedHosts := viper.GetString("ALLOWED_HOSTS")
	router := gin.New()
	router.SetTrustedProxies([]string{allowedHosts})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	RegisterRoutes(router) //routes register

	return router
}
