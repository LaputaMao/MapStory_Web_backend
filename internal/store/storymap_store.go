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

// List 返回 StoryMap 概览列表（分页和搜索）
// 注意：我们使用 Select 来只查询需要的字段，提高效率
func (s *StoryMapStore) List(page, pageSize int, title string) ([]model.StoryMap, int64, error) {
	var storyMaps []model.StoryMap
	var total int64

	query := s.db.Model(&model.StoryMap{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	// 先计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查询分页数据，只选择需要的字段
	offset := (page - 1) * pageSize
	err := query.Select("id", "title", "subtitle", "title_background", "created_at").
		Order("created_at desc"). // 按创建时间降序排序
		Offset(offset).
		Limit(pageSize).
		Find(&storyMaps).Error

	if err != nil {
		return nil, 0, err
	}

	return storyMaps, total, nil
}

// Delete 根据 ID 从数据库中删除 StoryMap
func (s *StoryMapStore) Delete(id uint) error {
	return s.db.Delete(&model.StoryMap{}, id).Error
}
