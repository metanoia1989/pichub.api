package services

import (
	"pichub.api/constants"
	"pichub.api/infra/database"
	"pichub.api/models"
	"pichub.api/pkg/utils"
)

type repositoryService struct{}

var RepositoryService = new(repositoryService)

func (s *repositoryService) AddRepository(userID int, repoName string, repoURL string, repoBranch string) (*models.Repository, error) {
	// 验证仓库
	_, _, err := GithubService.ValidateRepository(repoURL, "")
	if err != nil {
		return nil, err
	}

	// 创建仓库记录
	repository := &models.Repository{
		UserID:     userID,
		RepoName:   repoName,
		RepoURL:    repoURL,
		RepoBranch: utils.If(repoBranch == "", constants.DefaultRepoBranch, repoBranch),
	}

	if err := database.DB.Create(repository).Error; err != nil {
		return nil, err
	}

	return repository, nil
}

func (s *repositoryService) InitRepository(userID int, repoID int) (*models.Repository, error) {
	// 获取仓库信息
	repository, err := s.GetRepository(userID, repoID)
	if err != nil {
		return nil, err
	}

	// 初始化仓库数据
	if err := GithubService.InitializeRepository(repository, ""); err != nil {
		return nil, err
	}

	return repository, nil
}

func (s *repositoryService) ListRepositories(userID int) ([]models.Repository, error) {
	var repositories []models.Repository
	if err := database.DB.Where("user_id = ?", userID).Find(&repositories).Error; err != nil {
		return nil, err
	}
	return repositories, nil
}

func (s *repositoryService) GetRepository(userID int, repoID int) (*models.Repository, error) {
	var repository models.Repository
	if err := database.DB.Where("id = ? AND user_id = ?", repoID, userID).First(&repository).Error; err != nil {
		return nil, err
	}
	return &repository, nil
}

func (s *repositoryService) UpdateRepository(userID int, repoID int, repoURL string) error {
	// 验证仓库
	owner, repo, err := GithubService.ValidateRepository(repoURL, "")
	if err != nil {
		return err
	}

	// 更新仓库信息
	updates := map[string]interface{}{
		"repo_name": owner + "/" + repo,
		"repo_url":  repoURL,
	}

	result := database.DB.Model(&models.Repository{}).
		Where("id = ? AND user_id = ?", repoID, userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return database.DB.Error
	}

	return nil
}
