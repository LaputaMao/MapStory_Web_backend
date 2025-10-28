// Package service my-web-app/internal/service/storymap_service.go
package service

import (
	"StoryMap/internal/model"
	"StoryMap/internal/store"
)

type StoryMapService struct {
	store *store.StoryMapStore
}

func NewStoryMapService(s *store.StoryMapStore) *StoryMapService {
	return &StoryMapService{store: s}
}

// CreateStoryMap 将 StoryMap 模型存入数据库
func (s *StoryMapService) CreateStoryMap(storyMap *model.StoryMap) (*model.StoryMap, error) {
	if err := s.store.Create(storyMap); err != nil {
		return nil, err
	}
	return storyMap, nil
}

// GetStoryMapByID 根据 ID 获取 StoryMap
func (s *StoryMapService) GetStoryMapByID(id uint) (*model.StoryMap, error) {
	return s.store.GetByID(id)
}

// ListStoryMaps 获取 StoryMap 概览列表
func (s *StoryMapService) ListStoryMaps(page, pageSize int, title string) ([]model.StoryMap, int64, error) {
	return s.store.List(page, pageSize, title)
}

// DeleteStoryMapByID 根据 ID 删除 StoryMap
func (s *StoryMapService) DeleteStoryMapByID(id uint) error {
	return s.store.Delete(id)
}
