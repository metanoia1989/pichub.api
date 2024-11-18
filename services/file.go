package services

import (
	"fmt"
	_ "image/gif"  // 注册GIF格式
	_ "image/jpeg" // 注册JPEG格式
	_ "image/png"  // 注册PNG格式
	"io"
	"mime/multipart"
	"path/filepath"

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
	if err := GithubService.UploadFile(userID, repo.RepoURL, filePath, src); err != nil {
		return nil, err
	}

	// 创建文件记录
	fileRecord := &models.File{
		RepoID:      repoID,
		UserID:      userID,
		Filename:    filename,
		URL:         utils.BuildFileURL(repo.RepoURL, filePath),
		HashValue:   hashValue,
		RawFilename: file.Filename,
		Filesize:    uint(file.Size),
		Mime:        file.Header.Get("Content-Type"),
		Filetype:    fileType,
	}

	// 如果是图片，获取尺寸信息
	if fileType == 1 {
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
