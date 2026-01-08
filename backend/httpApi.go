package main

import (
	gin "github.com/gin-gonic/gin"
)

func SetupHttpRequestServer(r *gin.Engine) {
	SetupHttpWebsiteRequestServer(r)
	SetupHttpApiRequestServer(r)
}

func SetupHttpWebsiteRequestServer(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.File("./../frontend/dist/index.html")
	})
	r.Static("/assets", "./../frontend/dist/assets/")
}

func SetupHttpApiRequestServer(r *gin.Engine) {
	
}
