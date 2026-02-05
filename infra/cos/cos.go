package cos

import (
	"ai_interview/conf"
	"ai_interview/pkg/zlog"
	"bytes"
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

var cosClient *cos.Client

var (
	MaxFileSize = int64(50 * 1024 * 1024) // 50MB

	PdfMagicNumber  = []byte{0x25, 0x50, 0x44, 0x46, 0x2D}
	DocMagicNumber  = []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
	DocxMagicNumber = []byte{0x50, 0x4B, 0x03, 0x04}
)

func GetCosClient() *cos.Client {
	return cosClient
}

func InitCos() {
	u, err := url.Parse(conf.GetConfig().Cos.BucketUrl)
	if err != nil {
		panic("cos url解析失败" + err.Error())
	}

	b := &cos.BaseURL{BatchURL: u}

	cosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.GetConfig().Cos.SecretId,
			SecretKey: conf.GetConfig().Cos.SecretKey,
		},
	})

	if cosClient == nil {
		panic("cos客户端初始化失败")
	}
	zlog.Infof("cos客户端初始化成功")
}

func validateFileSize(size int64) bool {
	return size <= MaxFileSize
}

// 简历文件类型判断
func validateFileRealType(file io.ReadSeeker, filename string) (bool, error) {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))

	buf := make([]byte, 32) // 读取前32个字节
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("读取文件失败: %v", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, fmt.Errorf("文件指针复位失败: %v", err)
	}

	// 魔数校验
	if ext == ".pdf" {
		return bytes.HasPrefix(buf[:n], PdfMagicNumber), nil
	} else if ext == ".doc" {
		return bytes.HasPrefix(buf[:n], DocMagicNumber), nil
	} else if ext == ".docx" {
		return bytes.HasPrefix(buf[:n], DocxMagicNumber), nil
	} else {
		return false, fmt.Errorf("不支持的文件类型")
	}

}

// 返回公网url 和error
func UploadResume(file multipart.File, fileHeader *multipart.FileHeader, resumeID string) (string, error) {
	if !validateFileSize(fileHeader.Size) {
		return "", fmt.Errorf("文件大小超出限制")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	ok, err := validateFileRealType(file, fileHeader.Filename)

	if !ok {
		if err != nil {
			return "", err
		} else {
			return "", fmt.Errorf("不支持的文件类型")
		}
	}

	// 构建cos文件路径
	cosFilePath := fmt.Sprintf("%s/%s", ext, resumeID)

	_, err = cosClient.Object.Put(context.Background(), cosFilePath, file, nil)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %v", err)
	}

	return fmt.Sprintf("%s/%s", conf.GetConfig().Cos.BucketUrl, cosFilePath), nil

}
