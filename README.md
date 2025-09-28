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

### GORM (Go Object-Relational Mapper)
    
    它具体做了什么？
    建表 (CREATE TABLE)：
    如果数据库中不存在名为 images 的表（GORM 默认将结构体名 Image 转换为复数的小写形式），
    GORM 会生成并执行 CREATE TABLE 语句来创建它。它会根据 model.Image 结构体的字段名、类型和 GORM Tag 来定义表的列名、数据类型、主键、索引等。

    更新表结构 (ALTER TABLE)：
    如果 images 表已经存在，AutoMigrate 会检查您的 model.Image 结构体与现有表结构的差异，并执行 ALTER TABLE 语句：

    新增列：如果结构体中新增了字段，它会向表中添加新列。

    创建索引/外键：如果结构体中添加了相关的 GORM Tag，它会创建新的索引或外键。

    ⚠️ 一个重要限制：它不会删除或修改列