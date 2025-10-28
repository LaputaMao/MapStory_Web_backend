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
	"time"

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

// StoryMapOverviewResponse 定义了概览列表返回的结构
type StoryMapOverviewResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Subtitle        string    `json:"subtitle"`
	TitleBackground string    `json:"titleBackground"`
	CreatedAt       time.Time `json:"createdAt"`
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

	// 关键修改：只返回 content 字段
	// Gin 会自动处理 json.RawMessage，将其作为原始 JSON 输出
	c.Data(http.StatusOK, "application/json; charset=utf-8", storyMap.Content)
}

// List 获取 StoryMap 概览列表 (新增)
func (h *StoryMapHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	title := c.Query("title")

	storyMaps, total, err := h.service.ListStoryMaps(page, pageSize, title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get story map list"})
		return
	}

	// 将 model.StoryMap 转换为 StoryMapOverviewResponse
	overviews := make([]StoryMapOverviewResponse, len(storyMaps))
	for i, sm := range storyMaps {
		overviews[i] = StoryMapOverviewResponse{
			ID:              sm.ID,
			Title:           sm.Title,
			Subtitle:        sm.Subtitle,
			TitleBackground: sm.TitleBackground,
			CreatedAt:       sm.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": overviews,
		"meta": gin.H{
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Delete 根据 ID 删除 StoryMap (新增)
func (h *StoryMapHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid story map ID"})
		return
	}

	err = h.service.DeleteStoryMapByID(uint(id))
	if err != nil {
		// 在 GORM v2 中，删除一个不存在的记录不会返回错误，所以我们不需要检查 ErrRecordNotFound
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete story map"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Story map deleted successfully"})
}
