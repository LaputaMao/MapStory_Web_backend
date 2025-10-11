// Package router my-web-app/internal/router/router.go
package router

import (
	"StoryMap/internal/handler" // 导入你的 handler 包
	"github.com/gin-gonic/gin"
)

// SetupRouter 负责初始化和注册所有路由
// 注意，它接收所有需要的 handler 作为参数
func SetupRouter(uploadHandler *handler.UploadHandler, storyMapHandler *handler.StoryMapHandler, imageHandler *handler.ImageHandler, userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// 全局中间件 (如果需要的话)
	// r.Use(...)

	// 提供静态文件服务，让 /uploads/xxx.png 可以被访问
	r.Static("/uploads", "./uploads")

	// 设置一个更大的 multipart 内存限制
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 注册路由
	api := r.Group("/api/v1")
	{

		// 用户模块路由
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.GET("", userHandler.ListUsers) // GET /api/v1/users
		}

		// 图片上传路由
		api.POST("/uploadTest", imageHandler.Upload)

		// 独立的图片上传接口
		api.POST("/upload", uploadHandler.HandleImageUpload)
		api.DELETE("/upload", uploadHandler.HandleImageDelete) // <-- 图片删除路由

		// StoryMap 相关的接口
		storymaps := api.Group("/storymaps")
		{
			storymaps.POST("", storyMapHandler.Create)
			storymaps.GET("/:id", storyMapHandler.GetByID)
			// 未来可以添加 PUT (更新) 和 DELETE (删除)
		}

		// 假设未来我们有了一个 UserHandler
		// api.POST("/users", userHandler.Create)
		// api.GET("/users/:id", userHandler.Get)
	}

	return r
}
