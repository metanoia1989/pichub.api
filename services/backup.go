package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"pichub.api/infra/database"
	"pichub.api/models"
)

type BackupServiceImpl struct{}

var BackupService = &BackupServiceImpl{}

// CreateBackup 创建数据库备份
func (s *BackupServiceImpl) CreateBackup(repoID int) (*models.BackupRecord, error) {
	// 获取备份配置
	config := s.getBackupConfig()

	// 创建备份目录
	if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %v", err)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.sql", config.DatabaseName, timestamp)
	backupPath := filepath.Join(config.BackupDir, filename)

	// 执行mysqldump命令
	cmd := exec.Command("mysqldump",
		"-h", config.Host,
		"-P", config.Port,
		"-u", config.Username,
		fmt.Sprintf("-p%s", config.Password),
		"--databases", config.DatabaseName,
		"--single-transaction",
		"--quick",
		"--lock-tables=false",
		"-r", backupPath,
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("backup failed: %v", err)
	}

	// 压缩备份文件
	gzipPath := backupPath + ".gz"
	if err := s.compressFile(backupPath, gzipPath); err != nil {
		return nil, err
	}

	// 删除原始SQL文件
	os.Remove(backupPath)

	// 上传到GitHub
	githubPath := fmt.Sprintf("backups/%s/%s.gz", time.Now().Format("2006/01/02"), filename)
	if err := s.uploadToGitHub(repoID, gzipPath, githubPath); err != nil {
		return nil, err
	}

	// 创建备份记录
	record := &models.BackupRecord{
		RepoID:     repoID,
		BackupPath: githubPath,
		BackupDate: time.Now(),
	}

	if err := database.DB.Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to create backup record: %v", err)
	}

	// 清理本地备份文件
	os.Remove(gzipPath)

	return record, nil
}

// compressFile 压缩文件
func (s *BackupServiceImpl) compressFile(src, dst string) error {
	cmd := exec.Command("gzip", "-c", src)
	outfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outfile.Close()

	cmd.Stdout = outfile
	return cmd.Run()
}

// uploadToGitHub 上传文件到GitHub
func (s *BackupServiceImpl) uploadToGitHub(repoID int, localPath, remotePath string) error {
	// 获取仓库信息
	var repo models.Repository
	if err := database.DB.First(&repo, repoID).Error; err != nil {
		return fmt.Errorf("repository not found: %v", err)
	}

	// 读取文件内容
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 上传到GitHub
	return GithubService.UploadFile(repo.RepoURL, remotePath, file)
}

// getBackupConfig 获取备份配置
func (s *BackupServiceImpl) getBackupConfig() models.BackupConfig {
	return models.BackupConfig{
		BackupDir:    viper.GetString("BACKUP_DIR"),
		DatabaseName: viper.GetString("DB_DATABASE"),
		Username:     viper.GetString("DB_USERNAME"),
		Password:     viper.GetString("DB_PASSWORD"),
		Host:         viper.GetString("DB_HOST"),
		Port:         viper.GetString("DB_PORT"),
	}
}

// CleanOldBackups 清理旧的备份记录
func (s *BackupServiceImpl) CleanOldBackups(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return database.DB.Where("created_at < ?", cutoff).Delete(&models.BackupRecord{}).Error
}
