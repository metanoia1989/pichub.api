package services

import (
	"bytes"
	"fmt"
	_ "image/gif"  // 注册GIF格式
	_ "image/jpeg" // 注册JPEG格式
	_ "image/png"  // 注册PNG格式
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"pichub.api/infra/database"
	"pichub.api/models"
	"pichub.api/pkg/utils"
)

type FileServiceImpl struct{}

var FileService = &FileServiceImpl{}

// UploadFile 处理文件上传
func (s *FileServiceImpl) UploadFile(file *multipart.FileHeader, userID int, repoID int, isForce bool) (*models.File, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 读取文件内容用于计算哈希值和检测文件类型
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// 计算文件哈希值
	src.Seek(0, 0)
	hashValue, err := utils.CalculateGitHash(src, file.Size)
	if err != nil {
		return nil, err
	}

	// 重要：重置文件指针位置
	src.Seek(0, 0)

	// 如果不是强制上传，检查文件是否已存在
	if !isForce {
		var existingFile models.File
		if err := database.DB.Where("hash_value = ? AND repo_id = ?", hashValue, repoID).First(&existingFile).Error; err == nil {
			return &existingFile, nil
		}
	}
	// 已存在的，需要手动删除，程序不去处理了

	// 检测文件类型
	kind, _ := filetype.Match(buf[:n])
	fileType := utils.DetermineFileType(kind)

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" && kind != types.Unknown {
		ext = "." + kind.Extension
	}
	filename := fmt.Sprintf("%s%s", hashValue, ext)

	// 获取仓库信息
	var repo models.Repository
	if err := database.DB.First(&repo, repoID).Error; err != nil {
		return nil, fmt.Errorf("repository not found")
	}

	// 构建文件路径
	filePath := utils.BuildFilePath(filename)

	// 上传文件到GitHub
	src.Seek(0, 0)
	if err := GithubService.UploadFile(userID, repo.RepoURL, filePath, src); err != nil {
		return nil, err
	}

	repoPath := repo.GetRepositoryName()

	// 创建文件记录
	fileRecord := &models.File{
		RepoID:      repoID,
		UserID:      userID,
		Filename:    filename,
		URL:         filePath,
		RepoName:    repoPath,
		HashValue:   hashValue,
		RawFilename: file.Filename,
		Filesize:    uint(file.Size),
		Mime:        file.Header.Get("Content-Type"),
		Filetype:    fileType,
	}

	// 如果是图片，获取尺寸信息
	if fileType == 1 {
		src.Seek(0, 0) // 重置文件指针
		if width, height, err := utils.GetImageDimensions(src); err == nil {
			fileRecord.Width = uint(width)
			fileRecord.Height = uint(height)
		}
	}

	// 保存到数据库
	if err := database.DB.Create(fileRecord).Error; err != nil {
		return nil, err
	}

	return fileRecord, nil
}

// DeleteFile 删除文件
func (s *FileServiceImpl) DeleteFile(fileID, userID int) error {
	// 查找文件记录
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		return fmt.Errorf("file not found or no permission")
	}

	// 获取仓库信息
	var repo models.Repository
	if err := database.DB.First(&repo, file.RepoID).Error; err != nil {
		return fmt.Errorf("repository not found")
	}

	// 从GitHub删除文件
	err := GithubService.DeleteFile(userID, repo.RepoURL, file.URL)
	if err != nil {
		// 如果是404错误，直接继续删除数据库记录
		if !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("failed to delete file from GitHub: %v", err)
		}
	}

	// 删除数据库记录
	if err := database.DB.Delete(&file).Error; err != nil {
		return fmt.Errorf("failed to delete file record: %v", err)
	}

	return nil
}

// UploadStream 处理流式文件上传
func (s *FileServiceImpl) UploadStream(reader io.Reader, filename string, contentType string, fileSize int64, userID int, repoID int, isForce bool) (*models.File, error) {
	// 读取前512字节用于文件类型检测
	buf := make([]byte, 512)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// 创建一个新的reader，将已读取的内容和剩余内容组合
	combinedReader := io.MultiReader(bytes.NewReader(buf[:n]), reader)

	// 计算文件哈希值
	hashValue, err := utils.CalculateGitHash(combinedReader, fileSize)
	if err != nil {
		return nil, err
	}

	// 重新创建组合reader用于后续操作
	combinedReader = io.MultiReader(bytes.NewReader(buf[:n]), reader)

	// 如果不是强制上传，检查文件是否已存在
	if !isForce {
		var existingFile models.File
		if err := database.DB.Where("hash_value = ? AND repo_id = ?", hashValue, repoID).First(&existingFile).Error; err == nil {
			return &existingFile, nil
		}
	}

	// 检测文件类型
	kind, _ := filetype.Match(buf[:n])
	fileType := utils.DetermineFileType(kind)

	// 生成唯一文件名
	ext := filepath.Ext(filename)
	if ext == "" && kind != types.Unknown {
		ext = "." + kind.Extension
	}
	newFilename := fmt.Sprintf("%s%s", hashValue, ext)

	// 获取仓库信息
	var repo models.Repository
	if err := database.DB.First(&repo, repoID).Error; err != nil {
		return nil, fmt.Errorf("repository not found")
	}

	// 构建文件路径
	filePath := utils.BuildFilePath(newFilename)

	// 上传文件到GitHub
	if err := GithubService.UploadFile(userID, repo.RepoURL, filePath, combinedReader); err != nil {
		return nil, err
	}

	repoPath := repo.GetRepositoryName()

	// 创建文件记录
	fileRecord := &models.File{
		RepoID:      repoID,
		UserID:      userID,
		Filename:    newFilename,
		URL:         filePath,
		RepoName:    repoPath,
		HashValue:   hashValue,
		RawFilename: filename,
		Filesize:    uint(fileSize),
		Mime:        contentType,
		Filetype:    fileType,
	}

	// 如果是图片，获取尺寸信息
	if fileType == 1 {
		fullURL := fileRecord.GetFileURL(ConfigService.GetFileCDNHostname(0))

		// 重试最多3次
		for i := 0; i < 3; i++ {
			if width, height, err := utils.GetImageDimensionsFromURL(fullURL); err == nil {
				fileRecord.Width = uint(width)
				fileRecord.Height = uint(height)
				break
			}
			time.Sleep(time.Second * time.Duration(i+1)) // 递增延迟
		}
	}

	// 保存到数据库
	if err := database.DB.Create(fileRecord).Error; err != nil {
		return nil, err
	}

	return fileRecord, nil
}

// ListFiles 列出文件
func (s *FileServiceImpl) ListFiles(userID int, repoID int, page int, pageSize int) ([]models.File, int64, error) {
	var total int64
	var files []models.File

	// 构建基础查询
	query := database.DB.Where("user_id = ?", userID)

	// 如果指定了仓库ID，添加仓库筛选条件
	if repoID > 0 {
		query = query.Where("repo_id = ?", repoID)
	}

	// 获取总记录数
	if err := query.Model(&models.File{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加分页查询，并按ID降序排序
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}
