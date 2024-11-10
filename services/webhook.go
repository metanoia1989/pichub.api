package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"pichub.api/infra/database"
	"pichub.api/models"
)

type webhookService struct{}

var WebhookService = &webhookService{}

// ValidateSignature 验证 webhook 签名
func (s *webhookService) ValidateSignature(payload []byte, signature string) bool {
	secret := []byte(viper.GetString("GITHUB_WEBHOOK_SECRET"))

	// 计算 HMAC
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	// 比较签名
	return hmac.Equal([]byte(signature), []byte("sha256="+expectedMAC))
}

// HandlePush 处理 push 事件
func (s *webhookService) HandlePush(payload *models.WebhookPayload) error {
	// 只处理 master 分支的推送
	if !strings.HasSuffix(payload.Ref, "master") {
		return nil
	}

	// 获取仓库信息
	var repo models.Repository
	if err := database.DB.Where("repo_url LIKE ?", "%"+payload.Repository.FullName).First(&repo).Error; err != nil {
		return fmt.Errorf("repository not found: %v", err)
	}

	// 开启事务
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 处理所有提交
		for _, commit := range payload.Commits {
			// 处理新增文件
			if err := s.handleAddedFiles(tx, commit.Added, &repo); err != nil {
				return err
			}

			// 处理删除文件
			if err := s.handleRemovedFiles(tx, commit.Removed, &repo); err != nil {
				return err
			}

			// 处理修改文件
			if err := s.handleModifiedFiles(tx, commit.Modified, &repo); err != nil {
				return err
			}
		}
		return nil
	})
}

// handleAddedFiles 处理新增文件
func (s *webhookService) handleAddedFiles(tx *gorm.DB, files []string, repo *models.Repository) error {
	for _, filePath := range files {
		// 只处理 files 目录下的文件
		if !strings.HasPrefix(filePath, "files/") {
			continue
		}

		// 创建文件记录
		file := &models.File{
			RepoID:   repo.ID,
			UserID:   repo.UserID,
			Filename: getFilenameFromPath(filePath),
			URL:      buildGitHubRawURL(repo.RepoURL, filePath),
		}

		if err := tx.Create(file).Error; err != nil {
			return fmt.Errorf("failed to create file record: %v", err)
		}
	}
	return nil
}

// handleRemovedFiles 处理删除文件
func (s *webhookService) handleRemovedFiles(tx *gorm.DB, files []string, repo *models.Repository) error {
	for _, filePath := range files {
		if !strings.HasPrefix(filePath, "files/") {
			continue
		}

		if err := tx.Where("repo_id = ? AND filename = ?",
			repo.ID, getFilenameFromPath(filePath)).
			Delete(&models.File{}).Error; err != nil {
			return fmt.Errorf("failed to delete file record: %v", err)
		}
	}
	return nil
}

// handleModifiedFiles 处理修改文件
func (s *webhookService) handleModifiedFiles(tx *gorm.DB, files []string, repo *models.Repository) error {
	for _, filePath := range files {
		if !strings.HasPrefix(filePath, "files/") {
			continue
		}

		// 更新文件URL
		if err := tx.Model(&models.File{}).
			Where("repo_id = ? AND filename = ?", repo.ID, getFilenameFromPath(filePath)).
			Update("url", buildGitHubRawURL(repo.RepoURL, filePath)).Error; err != nil {
			return fmt.Errorf("failed to update file record: %v", err)
		}
	}
	return nil
}

// getFilenameFromPath 从文件路径中提取文件名
func getFilenameFromPath(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// buildGitHubRawURL 构建GitHub原始文件URL
func buildGitHubRawURL(repoURL, filePath string) string {
	return fmt.Sprintf("%s/raw/master/%s", repoURL, filePath)
}
