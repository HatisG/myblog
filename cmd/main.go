package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// post结构体，用于储存文章
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	//数据库连接
	dns := "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = sql.Open("mysql", dns)
	if err != nil {
		log.Fatal("数据库连接失败", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("数据库无法联通", err)
	}
	fmt.Println("数据库连接成功")

	//建表如果不存在
	_, err = db.Exec(`create table if not exists posts(
		id int(8) primary key auto_increment not null ,
		title varchar(255) not null,
		content text
	)`)
	if err != nil {
		log.Fatal("表创建失败", err)
	}

	//gin初始化
	gin.SetMode("release")
	r := gin.Default()

	//gin路由
	r.GET("/ping", pingHandler)
	r.GET("/post", postsHandler)

	//gin监听
	fmt.Println("server is running on :8080")
	r.Run(":8080")

}

// ping处理器函数
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "success"})
}

// post处理器函数
func postsHandler(c *gin.Context) {

	//用rows去接数据区posts里面的所有数据
	rows, err := db.Query("select id,title,content from posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	//posts切片
	var posts []Post

	//遍历rows，把数据拿出来
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, posts)

}
