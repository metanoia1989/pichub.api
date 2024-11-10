package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"pichub.api/infra/database"
	"pichub.api/models"
)

type fileService struct{}

var FileService = &fileService{}

// UploadFile 处理文件上传
func (s *fileService) UploadFile(file *multipart.FileHeader, userID int, repoID int, isForce bool) (*models.File, error) {
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
	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return nil, err
	}
	hashValue := hex.EncodeToString(hash.Sum(nil))

	// 如果不是强制上传，检查文件是否已存在
	if !isForce {
		var existingFile models.File
		if err := database.DB.Where("hash_value = ? AND repo_id = ?", hashValue, repoID).First(&existingFile).Error; err == nil {
			return &existingFile, nil
		}
	}

	// 检测文件类型
	kind, _ := filetype.Match(buf[:n])
	fileType := determineFileType(kind)

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
	filePath := buildFilePath(filename)

	// 上传文件到GitHub
	if err := GithubService.UploadFile(repo.RepoURL, filePath, src); err != nil {
		return nil, err
	}

	// 创建文件记录
	fileRecord := &models.File{
		RepoID:      repoID,
		UserID:      userID,
		Filename:    filename,
		URL:         buildFileURL(repo.RepoURL, filePath),
		HashValue:   hashValue,
		RawFilename: file.Filename,
		Filesize:    uint(file.Size),
		Mime:        file.Header.Get("Content-Type"),
		Filetype:    fileType,
	}

	// 如果是图片，获取尺寸信息
	if fileType == 1 {
		if width, height, err := getImageDimensions(src); err == nil {
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

// determineFileType 根据文件类型返回对应的类型代码
func determineFileType(kind types.Type) uint8 {
	if kind == types.Unknown {
		return 0
	}

	switch {
	case strings.HasPrefix(kind.MIME.Type, "image"):
		return 1
	case strings.HasPrefix(kind.MIME.Type, "video"):
		return 2
	case strings.HasPrefix(kind.MIME.Type, "audio"):
		return 3
	case strings.HasPrefix(kind.MIME.Type, "text"):
		return 4
	default:
		return 5
	}
}

// buildFilePath 构建文件存储路径
func buildFilePath(filename string) string {
	return fmt.Sprintf("files/%s/%s", filename[:2], filename)
}

// buildFileURL 构建文件访问URL
func buildFileURL(repoURL, filePath string) string {
	return fmt.Sprintf("%s/raw/master/%s", repoURL, filePath)
}
