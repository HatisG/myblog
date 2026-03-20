package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func About(c *gin.Context) {
	content, err := os.ReadFile("./web/static/about.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}
