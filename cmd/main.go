package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// post结构体，用于储存文章
type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	//数据库连接
	dsn := "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败", err)
	}
	fmt.Println("数据库连接成功")

	//gorm自动迁移
	if err := db.AutoMigrate(&Post{}); err != nil {
		log.Fatal("自动迁移失败", err)
	}
	fmt.Println("表迁移成功")

	//gin初始化
	gin.SetMode("release")
	r := gin.Default()

	//静态文件
	r.StaticFile("/post.html", "./web/post.html")

	//gin路由
	r.GET("/ping", pingHandler)
	r.GET("/post", postsHandler)
	r.POST("/post", createPostHandler)
	r.GET("/post/:id", getPostHandler)
	r.PUT("/post/:id", updatePostHandler)
	r.DELETE("/post/:id", deletePostHandler)

	//启动服务
	fmt.Println("server is running on :8080")
	r.Run(":8080")

}

// ping处理器函数
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "success"})
}

// post处理器函数
func postsHandler(c *gin.Context) {
	var posts []Post
	//gorm中Find方法直接遍历，并返回给posts这个临时结构体
	if err := db.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, posts)

}

// createPost处理函数
func createPostHandler(c *gin.Context) {

	//新文章的临时结构体
	var newPost Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//gorm使用Create方法创建新的Post
	if err := db.Create(&newPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPost)
}

// getPost的处理函数
func getPostHandler(c *gin.Context) {
	//从URL中获得id的参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	var post Post

	//使用First方法主键查询
	if err := db.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//返回文章json
	c.JSON(http.StatusOK, post)

}

// updatePost处理函数
func updatePostHandler(c *gin.Context) {
	//从URL上取id并转为整数
	isStr := c.Param("id")
	id, err := strconv.Atoi(isStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	//临时Post并绑定数据
	var updatePost Post
	if err := c.ShouldBindJSON(&updatePost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//gorm更新数据
	res := db.Model(&Post{}).Where("id = ?", id).Updates(updatePost)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	//返回res中的Rowsaffected，即此次修改影响的行数
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

// deletePost处理函数
func deletePostHandler(c *gin.Context) {
	//取id整数化
	isStr := c.Param("id")
	id, err := strconv.Atoi(isStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	//gorm通过id删除
	res := db.Delete(&Post{}, id)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	//查看删除行为对行数的影响
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
