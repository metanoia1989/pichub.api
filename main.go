package main

import (
	"time"

	"pichub.api/config"
	"pichub.api/infra/database"
	"pichub.api/infra/logger"
	"pichub.api/pkg/validator"
	"pichub.api/routers"
	"pichub.api/services"

	"github.com/spf13/viper"
)

func main() {

	//set timezone 设置时区
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Shanghai")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	// 加载配置，连接数据库
	config, err := config.SetupConfig()
	if err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	// 初始化各个服务
	if err := services.InitRedis(); err != nil {
		logger.Fatalf("redis connection error: %s", err)
	}

	if err := database.DbConnection(); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

	// 初始化验证器翻译器
	if err := validator.InitTrans(); err != nil {
		panic(err)
	}

	//later separate migration，不迁移，直接使用sql语句来操作表结构即可
	// migrations.Migrate()

	// 设置路由
	router := routers.SetupRoute()
	router.Static("/static", "./static")
	if err := router.Run(config.Server.ServerConfig()); err != nil {
		logger.Fatalf("Failed to start HTTP server: %v", err)
	}

	// 启动定时任务
	// services.SchedulerService.StartScheduler()
	// defer services.SchedulerService.StopScheduler()

}
