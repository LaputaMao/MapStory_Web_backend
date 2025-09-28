# MapStory

## 目录结构:

    StoryMap/
    ├── cmd/          # 项目的主入口 (main package)
    │   └── api/
    │       └── main.go
    ├── internal/     # 项目的内部私有代码，这是 Go 的一个特殊目录
    │   ├── handler/  # 处理 HTTP 请求的 handler (控制器)
    │   ├── model/    # 数据库模型 (struct)
    │   ├── service/  # 业务逻辑层
    │   └── store/    # 数据存储层 (数据库操作)
    ├── pkg/          # 可以被外部项目引用的公共代码 (初期可以不用)
    ├── configs/      # 配置文件 (如 config.yaml)
    ├── go.mod        # Go Modules 依赖文件
    └── go.sum        # 依赖包哈希文件