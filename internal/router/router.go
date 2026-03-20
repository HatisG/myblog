package router

import (
	"myblog/internal/handler"
	"myblog/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	gin.SetMode("release")
	r := gin.Default()

	// 静态文件
	r.LoadHTMLGlob("./web/*.html")
	r.Static("/static", "./web/static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/index.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/create-post.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create-post.html", nil)
	})
	r.GET("/search-post.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "search-post.html", nil)
	})
	r.GET("/about.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", nil)
	})

	// 初始化 service 和 handler
	postService := service.NewPostService(db)
	postHandler := handler.NewPostHandler(postService)

	// API 路由
	r.GET("/ping", postHandler.Ping)
	r.GET("/post", postHandler.List)
	r.POST("/post", postHandler.Create)
	r.GET("/post/:id", postHandler.Get)
	r.PUT("/post/:id", postHandler.Update)
	r.DELETE("/post/:id", postHandler.Delete)
	r.GET("/api/about", handler.About)

	// 打印所有路由
	for _, route := range r.Routes() {
		println(route.Method, route.Path)
	}

	return r
}
