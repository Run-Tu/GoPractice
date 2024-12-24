package main

import (
	"log"
	"project01/db"
	"project01/router"

	"github.com/spf13/viper"
)

func main() {
	r := router.SetupRouter()
	// 启动服务
	port := viper.GetString("app.port")
	if err := r.Run(":"+port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

