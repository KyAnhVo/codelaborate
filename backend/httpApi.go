package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/KyAnhVo/codelaborate/collab"
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
	api := r.Group("/api")
	{
		// Login endpoint
		api.POST("/login", func(c *gin.Context) {
			var loginReq struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			if err := c.BindJSON(&loginReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			// Simple authentication (in production, use proper auth)
			if loginReq.Username == "usrname" && loginReq.Password == "password" {
				// Generate session ID
				sessionID := generateSessionID()
				c.JSON(http.StatusOK, gin.H{
					"sessionID": sessionID,
					"username":  loginReq.Username,
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			}
		})

		// Create room endpoint
		api.POST("/rooms", func(c *gin.Context) {
			var createReq struct {
				SessionID string `json:"sessionID"`
			}

			if err := c.BindJSON(&createReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			// Generate room ID
			roomID := generateRoomID()
			c.JSON(http.StatusOK, gin.H{
				"roomID": roomID,
			})
		})

		// Get room endpoint
		api.GET("/rooms/:id", func(c *gin.Context) {
			roomID := c.Param("id")
			_, exists := collab.GlobalRoomManager.GetRoom(roomID)

			if exists {
				c.JSON(http.StatusOK, gin.H{
					"roomID": roomID,
					"exists": true,
				})
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Room not found",
				})
			}
		})

		// WebSocket endpoint for collaboration
		api.GET("/ws", collab.SetupCollabRoomServer)
	}
}

func generateSessionID() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func generateRoomID() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
