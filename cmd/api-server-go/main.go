package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

func SetupServer() *gin.Engine {
	r := gin.Default()
	r.GET("/", HelloHandler)
	return r
}

func main() {
	SetupServer().Run()
}
