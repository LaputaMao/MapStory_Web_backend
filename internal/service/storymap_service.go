// Package service my-web-app/internal/service/storymap_service.go
package service

import (
	"StoryMap/internal/model"
	"StoryMap/internal/store"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type StoryMapService struct {
	store   *store.StoryMapStore
	baseURL string
}

// NewStoryMapService 构造函数增加 baseURL 参数
func NewStoryMapService(s *store.StoryMapStore, baseURL string) *StoryMapService {
	return &StoryMapService{store: s, baseURL: baseURL}
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
//func (s *StoryMapService) DeleteStoryMapByID(id uint) error {
//	return s.store.Delete(id)

// DeleteStoryMapByID 根据 ID 删除 StoryMap 及其关联的图片文件 (修改后)
func (s *StoryMapService) DeleteStoryMapByID(id uint) error {
	// 1. 先从数据库获取 StoryMap 记录，以便访问其 content
	storyMap, err := s.store.GetByID(id)
	if err != nil {
		// 如果记录本就不存在，可以直接返回成功
		return nil
	}

	// 2. 解析 content JSON 并提取所有图片 URL
	urlsToDelete := s.extractImageURLsFromContent(storyMap.Content)

	// 3. 遍历 URL 列表并删除对应的文件
	for _, fileURL := range urlsToDelete {
		// 解析 URL 获取路径
		parsedURL, err := url.Parse(fileURL)
		if err != nil {
			// URL 格式错误，跳过
			continue
		}

		// *** 安全性检查 ***
		// 清理路径，防止路径遍历攻击 (e.g., /uploads/../../main.go)
		filePath := filepath.Clean(strings.TrimPrefix(parsedURL.Path, "/"))

		// 确保最终路径在 'uploads' 目录下
		if !strings.HasPrefix(filePath, "uploads"+string(os.PathSeparator)) {
			continue // 不是 uploads 目录下的文件，跳过
		}

		// 删除文件，忽略“文件不存在”的错误
		_ = os.Remove(filePath)
	}

	// 4. 最后，从数据库中删除 StoryMap 记录
	return s.store.Delete(id)
}

// extractImageURLsFromContent 递归地从 JSON 中提取所有指向本应用的图片 URL
func (s *StoryMapService) extractImageURLsFromContent(content json.RawMessage) []string {
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil
	}

	var urls []string
	s.findURLsInJSON(data, &urls)
	return urls
}

// findURLsInJSON 是一个递归辅助函数
func (s *StoryMapService) findURLsInJSON(data interface{}, urls *[]string) {
	switch v := data.(type) {
	case map[string]interface{}:
		for _, value := range v {
			s.findURLsInJSON(value, urls)
		}
	case []interface{}:
		for _, item := range v {
			s.findURLsInJSON(item, urls)
		}
	case string:
		// 检查字符串是否是本服务器的 uploads URL
		if strings.HasPrefix(v, s.baseURL+"/uploads/") {
			*urls = append(*urls, v)
		}
	}
}
