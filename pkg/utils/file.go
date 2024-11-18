package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"image"
	"io"
	"path/filepath"
	"strings"

	"mime/multipart"

	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
)

// DetermineFileType 根据文件类型返回对应的类型代码
func DetermineFileType(kind types.Type) uint8 {
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

func MimeToString(mime types.MIME) string {
	return fmt.Sprintf("%s/%s", mime.Type, mime.Subtype)
}

// GetFileType 根据文件名返回对应的 filetype.Type
func GetFileType(filename string) types.Type {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return types.Unknown
	}

	// 移除扩展名前的点号
	ext = strings.TrimPrefix(ext, ".")

	// 图片类型
	switch ext {
	case "jpg", "jpeg":
		return matchers.TypeJpeg
	case "png":
		return matchers.TypePng
	case "gif":
		return matchers.TypeGif
	case "webp":
		return matchers.TypeWebp
	case "avif":
		return types.Type{
			Extension: ext,
			MIME: types.MIME{
				Type:    "image",
				Subtype: "avif",
			},
		}
	case "svg":
		return types.Type{
			Extension: ext,
			MIME: types.MIME{
				Type:    "image",
				Subtype: "svg",
			},
		}

	// 视频类型
	case "mp4":
		return matchers.TypeMp4
	case "avi":
		return matchers.TypeAvi
	case "mov":
		return types.Type{
			Extension: ext,
			MIME: types.MIME{
				Type:    "video",
				Subtype: "quicktime",
			},
		}

	// 音频类型
	case "mp3":
		return matchers.TypeMp3
	case "wav":
		return matchers.TypeWav
	case "ogg":
		return matchers.TypeOgg
	case "m4a":
		return matchers.TypeM4a
	case "flac":
		return matchers.TypeFlac

	// 文档类型
	case "pdf":
		return matchers.TypePdf
	case "doc", "docx":
		return matchers.TypeDoc
	case "xls", "xlsx":
		return matchers.TypeXls

	// 压缩文件
	case "zip":
		return matchers.TypeZip
	case "rar":
		return matchers.TypeRar
	case "gz":
		return matchers.TypeGz
	case "tar":
		return matchers.TypeTar

	// 文本类型
	case "txt", "md", "html", "css", "js", "json", "xml", "yaml", "yml":
		return types.Type{
			Extension: ext,
			MIME: types.MIME{
				Type:    "text",
				Subtype: "plain",
			},
		}
	}

	// 如果没有匹配到，返回 Unknown 类型
	return types.Unknown
}

// BuildFilePath 构建文件存储路径
func BuildFilePath(filename string) string {
	return fmt.Sprintf("files/%s/%s", filename[:2], filename)
}

// BuildFileURL 构建文件访问URL
func BuildFileURL(repoURL, filePath string) string {
	return fmt.Sprintf("%s/raw/master/%s", repoURL, filePath)
}

// GetImageDimensions 获取图片尺寸
func GetImageDimensions(file multipart.File) (int, int, error) {
	// 重置文件指针到开始位置
	file.Seek(0, 0)

	// 解码图片
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

// CalculateGitHash 计算 Git 对象的 SHA1 散列值
func CalculateGitHash(content io.Reader, size int64) (string, error) {
	header := fmt.Sprintf("blob %d\x00", size)
	hash := sha1.New()

	// 写入 Git 对象头部
	if _, err := hash.Write([]byte(header)); err != nil {
		return "", err
	}

	// 写入文件内容
	if _, err := io.Copy(hash, content); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// FriendlyFileSize 将文件大小转换为人类友好的格式
// 输入bytes，根据大小自动转换为B、KB、MB、GB、TB等单位
func FriendlyFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	// 根据大小选择合适的单位
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.2f %s", float64(size)/float64(div), units[exp])
}
