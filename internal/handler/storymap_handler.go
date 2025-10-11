// Package handler my-web-app/internal/handler/storymap_handler.go
package handler

import (
	"StoryMap/internal/model"
	"StoryMap/internal/service"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// storyMapMeta 只用于从原始 JSON 中提取元数据
type storyMapMeta struct {
	Title           string `json:"title" binding:"required"`
	Subtitle        string `json:"subtitle"`
	TitleBackground string `json:"titleBackground"`
	//Content         json.RawMessage `json:"content" binding:"required"`
}

type StoryMapHandler struct {
	service *service.StoryMapService
}

func NewStoryMapHandler(s *service.StoryMapService) *StoryMapHandler {
	return &StoryMapHandler{service: s}
}

// Create 接收前端请求，创建 StoryMap
func (h *StoryMapHandler) Create(c *gin.Context) {
	// 1. 读取原始的请求体 (JSON)
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// 2. 从原始 JSON 中解析出元数据 (title, subtitle 等)
	var meta storyMapMeta
	if err := json.Unmarshal(rawBody, &meta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// 3. 验证必要字段
	if meta.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is a required field"})
		return
	}

	// 4. 将请求数据映射到数据库模型
	storyMapModel := &model.StoryMap{
		// 从元数据中填充独立字段
		Title:           meta.Title,
		Subtitle:        meta.Subtitle,
		TitleBackground: meta.TitleBackground,
		// 将完整的原始 JSON 保存到 content 字段
		Content: json.RawMessage(rawBody),
	}

	// 5. 调用 service 保存到数据库
	createdStoryMap, err := h.service.CreateStoryMap(storyMapModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create story map"})
		return
	}

	c.JSON(http.StatusCreated, createdStoryMap)
}

// GetByID 根据 ID 获取 StoryMap
func (h *StoryMapHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid story map ID"})
		return
	}

	storyMap, err := h.service.GetStoryMapByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Story map not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get story map"})
		return
	}

	c.JSON(http.StatusOK, storyMap)
}
