package services

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

type schedulerService struct {
	cron *cron.Cron
}

var SchedulerService = &schedulerService{
	cron: cron.New(cron.WithLocation(time.Local)),
}

// StartScheduler 启动定时任务调度器
func (s *schedulerService) StartScheduler() {
	// 添加数据库备份任务
	backupSchedule := viper.GetString("BACKUP_SCHEDULE")
	if backupSchedule == "" {
		backupSchedule = "0 0 * * *" // 默认每天凌晨执行
	}

	s.cron.AddFunc(backupSchedule, func() {
		log.Println("Starting database backup...")

		// 获取备份仓库ID
		backupRepoID := viper.GetInt("BACKUP_REPO_ID")
		if backupRepoID == 0 {
			log.Println("Backup repository not configured")
			return
		}

		// 执行备份
		record, err := BackupService.CreateBackup(backupRepoID)
		if err != nil {
			log.Printf("Backup failed: %v\n", err)
			return
		}

		log.Printf("Backup completed successfully: %s\n", record.BackupPath)

		// 清理30天前的备份
		if err := BackupService.CleanOldBackups(30); err != nil {
			log.Printf("Failed to clean old backups: %v\n", err)
		}
	})

	s.cron.Start()
}

// StopScheduler 停止定时任务调度器
func (s *schedulerService) StopScheduler() {
	if s.cron != nil {
		s.cron.Stop()
	}
}
