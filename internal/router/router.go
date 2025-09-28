// Package router my-web-app/internal/router/router.go
package router

import (
	"StoryMap/internal/handler" // 导入你的 handler 包
	"github.com/gin-gonic/gin"
)

// SetupRouter 负责初始化和注册所有路由
// 注意，它接收所有需要的 handler 作为参数
func SetupRouter(imageHandler *handler.ImageHandler) *gin.Engine {
	r := gin.Default()

	// 全局中间件 (如果需要的话)
	// r.Use(...)

	// 设置一个更大的 multipart 内存限制
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 注册路由
	api := r.Group("/api/v1")
	{
		// 图片上传路由
		api.POST("/upload", imageHandler.Upload)

		// 假设未来我们有了一个 UserHandler
		// api.POST("/users", userHandler.Create)
		// api.GET("/users/:id", userHandler.Get)
	}

	return r
}
