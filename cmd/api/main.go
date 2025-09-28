// my-web-app/cmd/api/main.go
package main

import "log"

func main() {
	// 创建一个新的服务器实例
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// 启动服务器
	server.Start()
}
