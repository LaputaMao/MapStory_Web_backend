# MapStory

## 目录结构:

    StoryMap/
    ├── cmd/          # 项目的主入口 (main package)
    │   └── api/
    │       ├── main.go
    │       └── server.go   # 依赖注入容器,简化main.go代码
    ├── internal/     # 项目的内部私有代码，这是 Go 的一个特殊目录
    │   ├── handler/  # 处理 HTTP 请求的 handler (控制器)
    │   ├── model/    # 数据库模型 (struct)
    │   ├── router/   # 设置所有API路由
    │   ├── service/  # 业务逻辑层
    │   └── store/    # 数据存储层 (数据库操作)
    ├── pkg/          # 可以被外部项目引用的公共代码 (初期可以不用)
    ├── configs/      # 配置文件 (如 config.yaml)
    ├── go.mod        # Go Modules 依赖文件
    └── go.sum        # 依赖包哈希文件

### GORM (Go Object-Relational Mapper)

    它具体做了什么？
    建表 (CREATE TABLE)：

    更新表结构 (ALTER TABLE)：

    ⚠️ 一个重要限制：它不会删除或修改列

### 更新日志

#### 2025.10.11 StoryMap创建模块 : json持久化 + 静态文件服务
    
    my-web-app/
    ├── cmd/
    │   └── api/
    │       ├── main.go
    │       └── server.go
    ├── configs/
    │   └── config.go
    ├── internal/
    │   ├── handler/
    │   │   ├── storymap_handler.go
    │   │   └── upload_handler.go  <-- 新增一个简单的上传处理器
    │   ├── model/
    │   │   └── storymap.go        <-- 修改后的模型
    │   ├── router/
    │   │   └── router.go          <-- 修改后的路由
    │   ├── service/
    │   │   └── storymap_service.go  <-- 修改后的服务
    │   └── store/
    │       └── storymap_store.go    <-- 修改后的存储
    └── uploads/                     <-- 图片会保存在这里
