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

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {

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

	_, err = db.Exec(`create table if not exists posts(
		id int(8) primary key auto_increment not null ,
		title varchar(255) not null,
		content text
	)`)
	if err != nil {
		log.Fatal("表创建失败", err)
	}

	gin.SetMode("release")
	r := gin.Default()
	r.GET("/ping", pingHandler)
	r.GET("/post", postsHandler)
	fmt.Println("server is running on :8080")
	r.Run(":8080")

}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "success"})
}
func postsHandler(c *gin.Context) {
	rows, err := db.Query("select id,title,content from posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var posts []Post
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
