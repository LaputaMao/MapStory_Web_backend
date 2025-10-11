// Package configs my-web-app/configs/config.go
package configs

// Config 是整个应用的配置
type Config struct {
	Server ServerConfig `json:"server"`
	MySQL  MySQLConfig  `json:"mysql"`
	App    AppConfig    `json:"app"`
}

// ServerConfig 是服务器相关的配置
type ServerConfig struct {
	Port string `json:"port"`
}

// MySQLConfig 是 MySQL 数据库的配置
type MySQLConfig struct {
	DSN string `json:"dsn"` // Data Source Name
}

// AppConfig 用于存放应用相关的配置
type AppConfig struct {
	BaseURL string `json:"base_url"` // 例如 "http://localhost:9090"
}

// LoadConfig 是一个简单的函数，用于加载配置
// 在真实项目中，这里会从文件或环境变量读取，我们先硬编码作为示例
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: ":9090", // 服务端口
		},
		MySQL: MySQLConfig{
			// DSN 格式: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
			// 请替换成你自己的 MySQL 信息
			DSN: "root:123456@tcp(127.0.0.1:3306)/map_story_data?charset=utf8mb4&parseTime=True&loc=Local",
		},
		App: AppConfig{
			BaseURL: "http://127.0.0.1:9090", //用于拼接图片生成静态文件URL
		},
	}
}
