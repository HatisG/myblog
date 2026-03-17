package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode("release")
	r := gin.Default()
	r.GET("/ping", pingHandler)
	r.GET("/post", postHandler)
	fmt.Println("server is running on :8080")
	r.Run(":8080")

}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "success"})
}
func postHandler(c *gin.Context) {
	c.JSON(200, []interface{}{})
}
