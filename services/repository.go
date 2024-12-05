package services

import (
	"errors"

	"pichub.api/constants"
	"pichub.api/infra/database"
	"pichub.api/models"
	"pichub.api/pkg/utils"
)

type repositoryService struct{}

var RepositoryService = new(repositoryService)

func (s *repositoryService) AddRepository(userID int, repoName string, repoURL string, repoBranch string) (*models.Repository, error) {

	// 检测记录是否已存在
	var repository *models.Repository
	if err := database.DB.Where("user_id = ? AND repo_url = ?", userID, repoURL).Find(repository).Error; err == nil {
		return repository, nil
	}

	// 先验证用户是否填写 github token
	token, err := ConfigService.GetGithubToken(userID)
	if err != nil || utils.IsEmpty(token) {
		return nil, errors.New("请先配置 github token")
	}

	// 验证仓库
	repoBranch = utils.If(repoBranch == "", constants.DefaultRepoBranch, repoBranch)
	_, _, err = GithubService.ValidateRepository(repoURL, token, repoBranch)
	if err != nil {
		return nil, err
	}

	// 创建仓库记录
	repository = &models.Repository{
		UserID:     userID,
		RepoName:   repoName,
		RepoURL:    repoURL,
		RepoBranch: repoBranch,
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

func (s *repositoryService) UpdateRepository(userID int, repoID int, repoURL string, repoBranch string) error {
	// 验证仓库
	token, err := ConfigService.GetGithubToken(userID)
	if err != nil {
		return err
	}

	repoBranch = utils.If(repoBranch == "", constants.DefaultRepoBranch, repoBranch)
	owner, repo, err := GithubService.ValidateRepository(repoURL, token, repoBranch)
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

func (s *repositoryService) DeleteRepository(userID int, repoID int) error {
	// 先删除文件
	if err := database.DB.Where("repo_id = ? and user_id = ?", repoID, userID).Delete(&models.File{}).Error; err != nil {
		return err
	}

	// 再删除仓库
	return database.DB.Where("id = ? AND user_id = ?", repoID, userID).Delete(&models.Repository{}).Error
}
