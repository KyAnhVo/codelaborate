package main

import (
	gin "github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default();

	// send html, css, and js files
	r.GET("/", func(c *gin.Context) {
		c.File("./../frontend/dist/index.html")
	})
	r.Static("/assets", "./../frontend/dist/assets/")

	// API-specific settings

	r.Run(":8080")
}
