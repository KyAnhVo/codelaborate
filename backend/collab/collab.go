package collab

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SetupCollabRoomServer(r *gin.Context) {
	conn, err := upgrader.Upgrade(r.Writer, r.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {

	}
}
