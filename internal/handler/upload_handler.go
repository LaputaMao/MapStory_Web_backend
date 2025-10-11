// Package handler my-web-app/internal/handler/upload_handler.go
package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadHandler struct {
	BaseURL string // 从配置中传入，例如 "http://localhost:9090"
}

// DeleteImageRequest 定义了删除图片接口的请求体
type DeleteImageRequest struct {
	URL string `json:"url" binding:"required,url"`
}

func NewUploadHandler(baseURL string) *UploadHandler {
	return &UploadHandler{BaseURL: baseURL}
}

// HandleImageUpload 只负责保存文件并返回 URL
func (h *UploadHandler) HandleImageUpload(c *gin.Context) {
	file, err := c.FormFile("image") // "image" 是前端约定的 key
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// 生成新文件名，防止冲突
	newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	// 定义保存目录并确保它存在
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 保存文件
	savePath := filepath.Join(uploadDir, newFileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 构造可访问的 URL
	// 注意：filepath.ToSlash 确保在 Windows 上也使用 / 作为路径分隔符
	fullURL := h.BaseURL + "/" + filepath.ToSlash(savePath)

	// 返回给前端，格式可以和前端约定，这里用一个简单的格式
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     fullURL,
	})
}

// HandleImageDelete 负责根据 URL 删除磁盘上的文件
func (h *UploadHandler) HandleImageDelete(c *gin.Context) {
	var req DeleteImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 1. 解析 URL，获取路径部分
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	// 2. *** 安全性检查 ***
	// 清理路径，防止路径遍历攻击 (e.g., /uploads/../../main.go)
	// filepath.Clean 会解析 ".."
	filePath := filepath.Clean(strings.TrimPrefix(parsedURL.Path, "/"))

	// 确保最终路径仍然在 'uploads' 目录下
	uploadDir := "uploads"
	if !strings.HasPrefix(filePath, uploadDir+string(os.PathSeparator)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: cannot delete files outside of the uploads directory"})
		return
	}

	// 3. 删除文件
	if err := os.Remove(filePath); err != nil {
		// 如果文件本身就不存在，也算成功，因为目标就是让它消失
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{"message": "File already deleted or does not exist"})
			return
		}
		// 其他错误（如权限问题）则报告服务器错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
