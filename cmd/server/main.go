package main

import (
	"fmt"
	"myblog/internal/config"
	"myblog/internal/model"
	"myblog/internal/router"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库
	db := config.InitDB(cfg)

	// 自动迁移
	if err := db.AutoMigrate(&model.Post{}, &model.Tag{}); err != nil {
		panic("自动迁移失败: " + err.Error())
	}
	fmt.Println("表迁移成功")

	// 设置路由
	r := router.Setup(db)

	// 启动服务
	fmt.Println("server is running on :8080")
	r.Run(":8080")
}
