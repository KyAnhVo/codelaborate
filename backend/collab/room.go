package collab

import (
	"log"
	"sync"
)

type Room struct {
	// ID of room (multiple rooms possible)
	Id string

	// all msgs sent into room goes through this
	MsgInChannel chan *UpdateMsg

	// this manages client, send update to clients, etc.
	Clients map[string]*Client
	mu      sync.RWMutex

	// the core text string
	Text FileText
}

func CreateRoom(roomId string) *Room {
	room := new(Room)
	room.Id = roomId
	room.MsgInChannel = make(chan *UpdateMsg, 128)
	room.Clients = make(map[string]*Client)
	room.Text = FileText{}
	return room
}

func (r *Room) Start() {
	for {
		msg := <-r.MsgInChannel
		r.processMsg(msg)
	}
}

func (r *Room) processMsg(msg *UpdateMsg) {
	if msg.Op == "exit" {
		r.removeClient(msg.SessionID)
		return
	}

	// Broadcast to all clients except sender
	r.mu.RLock()
	defer r.mu.RUnlock()

	for sessionID, client := range r.Clients {
		if sessionID != string(rune(msg.SessionID)) {
			select {
			case client.OutChannel <- msg:
			default:
				log.Printf("Failed to send message to client %s", sessionID)
			}
		}
	}
}

func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clients[client.SessionID] = client
	log.Printf("Client %s joined room %s. Total clients: %d", client.SessionID, r.Id, len(r.Clients))
}

func (r *Room) removeClient(sessionID int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Clients, string(rune(sessionID)))
	log.Printf("Client removed from room %s. Total clients: %d", r.Id, len(r.Clients))
}

func (r *Room) RemoveClientBySessionID(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Clients, sessionID)
	log.Printf("Client %s removed from room %s. Total clients: %d", sessionID, r.Id, len(r.Clients))
}
