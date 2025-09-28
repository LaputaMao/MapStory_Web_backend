// Package store my-web-app/internal/store/image_store.go
package store

import (
	"StoryMap/internal/model" // 导入我们定义的模型
	"gorm.io/gorm"
)

// ImageStore 定义了图片数据存储需要实现的方法
type ImageStore struct {
	db *gorm.DB
}

// NewImageStore 创建一个 ImageStore 实例
func NewImageStore(db *gorm.DB) *ImageStore {
	return &ImageStore{db: db}
}

// Create 在数据库中创建一条图片记录
func (s *ImageStore) Create(image *model.Image) error {
	// 使用 GORM 的 Create 方法来插入数据
	return s.db.Create(image).Error
}
