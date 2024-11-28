package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v65/github"
	"pichub.api/config"
	"pichub.api/infra/database"
	"pichub.api/infra/logger"
	"pichub.api/models"
	"pichub.api/pkg/utils"
)

type GithubServiceImpl struct {
	client *github.Client
}

var GithubService = &GithubServiceImpl{}

// 方法1：创建一个调试用的 Transport
type debugTransport struct {
	t http.RoundTripper
}

func (d *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 打印请求信息
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logger.Errorf("Failed to dump request: %v", err)
	} else {
		fmt.Printf("GitHub Request:\n%s\n", string(reqDump))
	}

	// 执行请求
	resp, err := d.t.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// 打印响应信息
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		logger.Errorf("Failed to dump response: %v", err)
	} else {
		fmt.Printf("GitHub Response:\n%s\n", string(respDump))
	}

	return resp, nil
}

func (s *GithubServiceImpl) getClient(token string) *github.Client {
	logger.Infof("is github debug: %v", config.Config.Server.GithubDebug)
	httpClient := &http.Client{}
	if config.Config.Server.GithubDebug {
		httpClient = &http.Client{
			Transport: &debugTransport{
				t: http.DefaultTransport,
			},
		}
	}

	client := github.NewClient(httpClient)
	if token != "" {
		client = client.WithAuthToken(token)
	}
	return client
}

// ValidateRepository 验证仓库是否存在且可访问
func (s *GithubServiceImpl) ValidateRepository(repoURL string, token string, branch string) (owner, repo string, err error) {
	// 从URL中提取owner和repo名称
	parts := strings.Split(strings.TrimSuffix(repoURL, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid repository URL")
	}

	owner = parts[len(parts)-2]
	repo = parts[len(parts)-1]

	client := s.getClient(token)
	ctx := context.Background()

	logger.Infof("ValidateRepository %s %s", owner, repo)

	// 检查仓库是否存在
	_, _, err = client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", "", fmt.Errorf("repository not found or not accessible")
	}

	// 检查分支是否存在
	if branch != "" {
		_, _, err = client.Repositories.GetBranch(ctx, owner, repo, branch, 3)
		if err != nil {
			return "", "", fmt.Errorf("branch '%s' not found or not accessible: %v", branch, err)
		}
	}

	return owner, repo, nil
}

// InitializeRepository 初始化仓库数据
func (s *GithubServiceImpl) InitializeRepository(repo *models.Repository, token string) error {
	parts := strings.Split(strings.TrimSuffix(repo.RepoURL, "/"), "/")
	owner := parts[len(parts)-2]
	repoName := parts[len(parts)-1]

	client := s.getClient(token)
	ctx := context.Background()

	// 获取仓库中的所有文件
	opts := &github.RepositoryContentGetOptions{}
	_, contents, _, err := client.Repositories.GetContents(ctx, owner, repoName, "/", opts)
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
func (s *GithubServiceImpl) processContent(ctx context.Context, client *github.Client, owner, repoName string, content *github.RepositoryContent, repo *models.Repository) error {
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
		// 提取相对路径：跳过前四个部分（host/owner/repo/branch）
		urlParts := strings.Split(*content.DownloadURL, "/")
		relativePath := strings.Join(urlParts[6:], "/")

		fileType := utils.GetFileType(*content.Name)
		mime := utils.MimeToString(fileType.MIME)
		filetypeInt := utils.DetermineFileType(fileType)

		file := &models.File{
			RepoID:      repo.ID,
			UserID:      repo.UserID,
			RepoName:    repo.GetRepositoryName(),
			Filename:    *content.Name,
			URL:         relativePath,
			RawFilename: *content.Name,
			HashValue:   *content.SHA,
			Filesize:    uint(*content.Size),
			Filetype:    filetypeInt,
			Mime:        mime,
		}

		// 检测是否已存在 repoID, userID, filename 相同的文件，如果存在，则更新
		var existingFile models.File
		if err := database.DB.Where("repo_id = ? AND user_id = ? AND filename = ?", repo.ID, repo.UserID, file.Filename).First(&existingFile).Error; err == nil {
			existingFile.HashValue = file.HashValue
			existingFile.Filesize = file.Filesize
			existingFile.Filetype = file.Filetype
			existingFile.RepoName = file.RepoName
			existingFile.Mime = file.Mime
			existingFile.Width = 0
			existingFile.Height = 0

			if err := database.DB.Save(&existingFile).Error; err != nil {
				return fmt.Errorf("failed to update file record: %v", err)
			}
			return nil
		}

		// 保存文件记录
		if err := database.DB.Create(file).Error; err != nil {
			return fmt.Errorf("failed to save file record: %v", err)
		}
	}
	return nil
}

// UploadFile 上传文件到GitHub仓库
func (s *GithubServiceImpl) UploadFile(userID int, repoURL string, remotePath string, file io.Reader) error {
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

	// 获取token
	token, err := ConfigService.GetGithubToken(userID)
	if err != nil {
		return fmt.Errorf("failed to get github token: %v", err)
	}

	client := s.getClient(token)
	ctx := context.Background()

	// 上传文件
	_, _, err = client.Repositories.CreateFile(ctx, owner, repo, remotePath, opts)
	if err != nil {
		return fmt.Errorf("failed to upload file to GitHub: %v", err)
	}

	return nil
}

// ValidateToken 验证 GitHub token 是否有效
func (s *GithubServiceImpl) ValidateToken(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("token is empty")
	}

	client := s.getClient(token)
	ctx := context.Background()

	// 尝试获取用户信息来验证 token
	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		if resp != nil && resp.StatusCode == 401 {
			return false, fmt.Errorf("invalid token")
		}
		return false, fmt.Errorf("failed to validate token: %v", err)
	}

	// 确保获取到用户信息
	if user == nil || user.Login == nil {
		return false, fmt.Errorf("failed to get user info")
	}

	return true, nil
}

func (s *GithubServiceImpl) GetToken(userID int) (string, error) {
	token, err := ConfigService.Get("github", "token", userID)
	if err != nil || utils.IsEmpty(token) {
		return "", fmt.Errorf("请先配置 github token")
	}
	return token.(string), nil
}

// DeleteFile 从GitHub仓库删除文件
func (s *GithubServiceImpl) DeleteFile(userID int, repoURL string, remotePath string) error {
	// 从URL中提取owner和repo名称
	parts := strings.Split(strings.TrimSuffix(repoURL, "/"), "/")
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	// 获取token
	token, err := s.GetToken(userID)
	if err != nil {
		return fmt.Errorf("failed to get github token: %v", err)
	}

	client := s.getClient(token)
	ctx := context.Background()

	// 获取文件的当前SHA
	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, remotePath, nil)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// 准备删除文件的参数
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("Delete file: %s", filepath.Base(remotePath))),
		SHA:     github.String(*content.SHA),
	}

	// 删除文件
	_, _, err = client.Repositories.DeleteFile(ctx, owner, repo, remotePath, opts)
	if err != nil {
		return fmt.Errorf("failed to delete file from GitHub: %v", err)
	}

	return nil
}
