// my-web-app/cmd/api/server.go
package main

import (
	"StoryMap/configs"
	"StoryMap/internal/handler"
	"StoryMap/internal/model"
	"StoryMap/internal/router" // 引入我们新的 router 包
	"StoryMap/internal/service"
	"StoryMap/internal/store"
	"github.com/gin-gonic/gin"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Server 是我们应用的核心结构体，它持有所有依赖
type Server struct {
	config *configs.Config
	engine *gin.Engine
}

// NewServer 创建一个新的 Server 实例
func NewServer() (*Server, error) {
	// 加载配置
	cfg := configs.LoadConfig()

	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	db.AutoMigrate(&model.Image{})

	// =================================================================
	// 依赖注入的核心部分
	// =================================================================
	// 1. 初始化 store 层
	imageStore := store.NewImageStore(db)
	// userStore := store.NewUserStore(db) // 未来...

	// 2. 初始化 service 层
	imageService := service.NewImageService(imageStore)
	// userService := service.NewUserService(userStore) // 未来...

	// 3. 初始化 handler 层
	imageHandler := handler.NewImageHandler(imageService)
	// userHandler := handler.NewUserHandler(userService) // 未来...

	// 4. 初始化路由
	// 把所有 handler 传递给路由设置函数
	r := router.SetupRouter(imageHandler)

	// 创建 Server 实例
	server := &Server{
		config: cfg,
		engine: r,
	}

	return server, nil
}

// Start 启动 HTTP 服务器
func (s *Server) Start() {
	log.Printf("Server is running at %s", s.config.Server.Port)
	if err := s.engine.Run(s.config.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
