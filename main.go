package main

import (
	"time"

	"pichub.api/config"
	"pichub.api/infra/database"
	"pichub.api/infra/logger"
	"pichub.api/routers"

	"github.com/spf13/viper"
)

func main() {

	//set timezone 设置时区
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Shanghai")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	// 加载配置，连接数据库
	_, err := config.SetupConfig()
	if err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	masterDSN := config.DbConfiguration()

	if err := database.DbConnection(masterDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

	//later separate migration，不迁移，直接使用sql语句来操作表结构即可
	// migrations.Migrate()

	// 设置路由
	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.ServerConfig()))

	// 启动定时任务
	// services.SchedulerService.StartScheduler()
	// defer services.SchedulerService.StopScheduler()

}
