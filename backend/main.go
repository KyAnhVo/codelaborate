package main

import (
	gin "github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default();
	SetupHttpRequestServer(r)
	
	r.Run(":8080")
}

