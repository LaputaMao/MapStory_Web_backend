// Package store my-web-app/internal/store/storymap_store.go
package store

import (
	"StoryMap/internal/model"
	"gorm.io/gorm"
)

type StoryMapStore struct {
	db *gorm.DB
}

func NewStoryMapStore(db *gorm.DB) *StoryMapStore {
	return &StoryMapStore{db: db}
}

// Create 在数据库中创建一条 StoryMap 记录
func (s *StoryMapStore) Create(storyMap *model.StoryMap) error {
	return s.db.Create(storyMap).Error
}

// GetByID 根据 ID 从数据库中查找 StoryMap
func (s *StoryMapStore) GetByID(id uint) (*model.StoryMap, error) {
	var storyMap model.StoryMap
	if err := s.db.First(&storyMap, id).Error; err != nil {
		return nil, err
	}
	return &storyMap, nil
}
