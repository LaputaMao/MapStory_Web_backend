// Package handler my-web-app/internal/handler/image_handler.go
package handler

import (
	"StoryMap/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ImageHandler 负责处理图片相关的 HTTP 请求
type ImageHandler struct {
	service *service.ImageService
}

// NewImageHandler 创建一个 ImageHandler 实例
func NewImageHandler(s *service.ImageService) *ImageHandler {
	return &ImageHandler{service: s}
}

// Upload 是处理图片上传的 Gin Handler
func (h *ImageHandler) Upload(c *gin.Context) {
	// 1. 从请求中获取上传的文件
	// "image" 是前端 <input type="file" name="image"> 中的 name 属性
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// 2. 调用 service 层处理业务逻辑
	image, err := h.service.UploadImage(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"data":    image,
	})
}
