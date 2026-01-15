package collab

import (
	"log"
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
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Read initial join message
	var joinMsg JoinMsg
	err = conn.ReadJSON(&joinMsg)
	if err != nil {
		log.Printf("Failed to read join message: %v", err)
		return
	}

	// Get or create room
	var room *Room
	if joinMsg.Op == "create" {
		room = GlobalRoomManager.GetOrCreateRoom(joinMsg.RoomID)
		log.Printf("Room %s created/retrieved", joinMsg.RoomID)
	} else if joinMsg.Op == "join" {
		roomID := r.Query("roomID")
		if roomID == "" {
			log.Printf("Room ID not provided for join operation")
			return
		}
		var exists bool
		room, exists = GlobalRoomManager.GetRoom(roomID)
		if !exists {
			log.Printf("Room %s does not exist", roomID)
			return
		}
		log.Printf("Client joining room %s", roomID)
	} else {
		log.Printf("Invalid operation: %s", joinMsg.Op)
		return
	}

	// Create client and add to room
	client := &Client{
		SessionID:  string(rune(joinMsg.SessionID)),
		OutChannel: make(chan *UpdateMsg, 128),
	}
	room.AddClient(client)
	defer room.RemoveClientBySessionID(client.SessionID)

	// Start goroutine to write messages to client
	go func() {
		for msg := range client.OutChannel {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Failed to write message to client: %v", err)
				return
			}
		}
	}()

	// Read messages from client and send to room
	for {
		var updateMsg UpdateMsg
		err := conn.ReadJSON(&updateMsg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Send message to room for broadcasting
		select {
		case room.MsgInChannel <- &updateMsg:
		default:
			log.Printf("Room message channel full")
		}
	}
}
