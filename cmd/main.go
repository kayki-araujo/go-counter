package main

import (
	"counter/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("internal/templates/*")

	router.GET("/", handlers.HomePage)
	router.POST("/inc", handlers.IncrementCounter)
	router.GET("/sse", handlers.SSEHandler)

	router.Run()
}
