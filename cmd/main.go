package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	dsn := "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = sql.Open("mysql", dsn)
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
	r.POST("/post", createPostHandler)
	r.GET("/post/:id", getPostHandler)
	r.PUT("/post/:id", updatePostHandler)
	r.DELETE("/post/:id", deletePostHandler)

	//gin监听
	fmt.Println("server is running on :8080")
	r.StaticFile("/post.html", "./web/post.html")
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

	c.JSON(200, posts)

}

// createPost处理函数
func createPostHandler(c *gin.Context) {

	//新文章的临时结构体
	var newPost struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//插入新文章，用result去接
	result, err := db.Exec("insert into posts (title,content) values (?,?)", newPost.Title, newPost.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//返回result的自增id
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//成功相应
	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"title":   newPost.Title,
		"content": newPost.Content,
	})

}

// getPost的处理函数
func getPostHandler(c *gin.Context) {
	//从URL中获得id的参数
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少文章ID"})
		return
	}

	var id int
	var err error
	//id从字符串转整数
	id, err = strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var post Post
	//根据id从database中拿post的id，title，content
	row := db.QueryRow("select id,title,content from posts where id = ?", id)
	//查询结果传到结构体内
	err = row.Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		if err == sql.ErrNoRows {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//临时Post
	var updatePost struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updatePost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Sql语句更新Post
	res, err := db.Exec("update posts set title=?,content=? where id=?", updatePost.Title, updatePost.Content, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//返回res中的rowsaffected，即此次修改影响的行数
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//SQL通过id删除
	res, err := db.Exec("delete from posts where id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//查看删除行为对行数的影响
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
