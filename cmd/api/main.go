package main

import "github.com/gin-gonic/gin"

func main() {
	// 1. 创建一个默认的 Gin 引擎
	r := gin.Default()

	// 2. 定义一个路由和处理函数
	// 当浏览器访问 http://localhost:8080/ping 时，会执行后面的函数
	r.GET("/ping", func(c *gin.Context) {
		// 以 JSON 格式返回一个消息
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 3. 启动 HTTP 服务器，默认监听在 8080 端口
	// 修改监听端口 , 形似 gogo hhh
	r.Run(":9090")
}
