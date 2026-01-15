package collab

import (
	"sync"
)

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

var GlobalRoomManager = &RoomManager{
	rooms: make(map[string]*Room),
}

func (rm *RoomManager) GetOrCreateRoom(roomID string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		room = CreateRoom(roomID)
		rm.rooms[roomID] = room
		go room.Start()
	}
	return room
}

func (rm *RoomManager) GetRoom(roomID string) (*Room, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	room, exists := rm.rooms[roomID]
	return room, exists
}

func (rm *RoomManager) DeleteRoom(roomID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.rooms, roomID)
}
