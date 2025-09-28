// my-web-app/cmd/api/main.go
package main

import (
	"StoryMap/configs"
	"StoryMap/internal/handler"
	"StoryMap/internal/model"
	"StoryMap/internal/service"
	"StoryMap/internal/store"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. 加载配置
	cfg := configs.LoadConfig()

	// 2. 初始化数据库连接 (GORM)
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移：自动创建或更新数据库表结构
	// 这会让 GORM 根据你的 model.Image 结构体自动创建 `images` 表
	db.AutoMigrate(&model.Image{})

	// 3. 依赖注入：将依赖关系串联起来
	// 这是核心！我们从底层(store)开始，一步步向上构建
	imageStore := store.NewImageStore(db)
	imageService := service.NewImageService(imageStore)
	imageHandler := handler.NewImageHandler(imageService)

	// 4. 初始化 Gin 引擎
	r := gin.Default()
	// 设置一个更大的 multipart 内存限制 (默认是 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 5. 注册路由
	// 创建一个路由组 /api/v1
	api := r.Group("/api/v1")
	{
		// 图片上传路由: POST /api/v1/upload
		api.POST("/upload", imageHandler.Upload)
	}

	// 6. 启动服务器
	log.Printf("Server is running at %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
