package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v65/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"pichub.api/infra/database"
	"pichub.api/models"
)

type githubService struct {
	client *github.Client
}

var GithubService = &githubService{}

func (s *githubService) getClient(token string) *github.Client {
	if token == "" {
		token = viper.GetString("GITHUB_TOKEN")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// ValidateRepository 验证仓库是否存在且可访问
func (s *githubService) ValidateRepository(repoURL string, token string) (owner, repo string, err error) {
	// 从URL中提取owner和repo名称
	parts := strings.Split(strings.TrimSuffix(repoURL, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid repository URL")
	}

	owner = parts[len(parts)-2]
	repo = parts[len(parts)-1]

	client := s.getClient(token)
	ctx := context.Background()

	// 检查仓库是否存在
	_, _, err = client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", "", fmt.Errorf("repository not found or not accessible")
	}

	return owner, repo, nil
}

// InitializeRepository 初始化仓库数据
func (s *githubService) InitializeRepository(repo *models.Repository, token string) error {
	parts := strings.Split(strings.TrimSuffix(repo.RepoURL, "/"), "/")
	owner := parts[len(parts)-2]
	repoName := parts[len(parts)-1]

	client := s.getClient(token)
	ctx := context.Background()

	// 获取仓库中的所有文件
	opts := &github.RepositoryContentGetOptions{}
	_, contents, _, err := client.Repositories.GetContents(ctx, owner, repoName, "", opts)
	if err != nil {
		return fmt.Errorf("failed to get repository contents: %v", err)
	}

	// 递归处理所有文件
	for _, content := range contents {
		if err := s.processContent(ctx, client, owner, repoName, content, repo); err != nil {
			return err
		}
	}

	return nil
}

// processContent 递归处理仓库内容
func (s *githubService) processContent(ctx context.Context, client *github.Client, owner, repoName string, content *github.RepositoryContent, repo *models.Repository) error {
	if *content.Type == "dir" {
		// 如果是目录，递归处理
		opts := &github.RepositoryContentGetOptions{}
		_, contents, _, err := client.Repositories.GetContents(ctx, owner, repoName, *content.Path, opts)
		if err != nil {
			return err
		}
		for _, c := range contents {
			if err := s.processContent(ctx, client, owner, repoName, c, repo); err != nil {
				return err
			}
		}
	} else {
		// 如果是文件，创建文件记录
		file := &models.File{
			RepoID:      repo.ID,
			UserID:      repo.UserID,
			Filename:    *content.Name,
			URL:         *content.DownloadURL,
			RawFilename: *content.Name,
		}

		// 保存文件记录
		if err := database.DB.Create(file).Error; err != nil {
			return fmt.Errorf("failed to save file record: %v", err)
		}
	}
	return nil
}

// UploadFile 上传文件到GitHub仓库
func (s *githubService) UploadFile(repoURL string, remotePath string, file io.Reader) error {
	// 从URL中提取owner和repo名称
	parts := strings.Split(strings.TrimSuffix(repoURL, "/"), "/")
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// 准备文件上传参数
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("Upload backup file: %s", filepath.Base(remotePath))),
		Content: content,
		// Branch:  github.String("main"), // 指定上传分支，可忽略
	}

	client := s.getClient("")
	ctx := context.Background()

	// 上传文件
	_, _, err = client.Repositories.CreateFile(ctx, owner, repo, remotePath, opts)
	if err != nil {
		return fmt.Errorf("failed to upload file to GitHub: %v", err)
	}

	return nil
}
