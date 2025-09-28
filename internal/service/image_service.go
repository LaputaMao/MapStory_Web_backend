// Package service my-web-app/internal/service/image_service.go
package service

import (
	"StoryMap/internal/model"
	"StoryMap/internal/store"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// ImageService 封装了图片相关的业务逻辑
type ImageService struct {
	store *store.ImageStore
}

// NewImageService 创建一个 ImageService 实例
func NewImageService(s *store.ImageStore) *ImageService {
	return &ImageService{store: s}
}

// UploadImage 处理图片上传的核心逻辑
func (s *ImageService) UploadImage(file *multipart.FileHeader) (*model.Image, error) {
	// 1. 生成新的文件名，防止重名
	// 规则：时间戳_原文件名
	newFileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)

	// 2. 定义保存路径 (例如: uploads/2023-10-27/xxxx.jpg)
	// 为了简单起见，我们直接保存在项目根目录下的 "uploads" 文件夹
	uploadDir := "uploads"
	// 确保目录存在
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}
	savePath := filepath.Join(uploadDir, newFileName)

	// 3. 保存文件到服务器本地
	// c.SaveUploadedFile 是 Gin 提供的便捷方法，我们在这里手动实现，更清晰
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// 将上传的文件内容拷贝到新创建的文件中
	if _, err = src.Seek(0, 0); err != nil { // 重置文件指针
		return nil, err
	}
	if _, err = dst.ReadFrom(src); err != nil {
		return nil, err
	}

	// 4. 创建模型对象，准备存入数据库
	imageModel := &model.Image{
		Name:   file.Filename,
		Path:   savePath,
		Format: filepath.Ext(file.Filename), // 获取文件扩展名
	}

	// 5. 调用 store 层，将图片信息存入数据库
	if err := s.store.Create(imageModel); err != nil {
		// 如果数据库保存失败，最好把刚刚保存的文件也删掉，以保持数据一致性
		os.Remove(savePath)
		return nil, err
	}

	return imageModel, nil
}
