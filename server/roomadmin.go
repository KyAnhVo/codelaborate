package main

import (
	"time"
	"sync"
	"errors"
)
// Only 128 rooms are available at a time max.
const maxRoomCount uint32 = 128

// Duration room is guranteed to exist before being available to be kicked out
const maxGuranteedRoomDuration = 30 * time.Minute

var roomManagers []*RoomManager = nil
var nextRoomID	uint32 	= 0
var firstRoomID	uint32 	= 0
var noRoom 		bool 	= true
var roomsLock sync.RWMutex

// CreateRoomState creates state of rooms
func CreateRoomState() {
	roomsLock.Lock()
	defer roomsLock.Unlock()
	roomManagers = make([]*RoomManager, maxRoomCount)
}

// JoinRoom attempts to add user to new room. Thread-safe.
func JoinRoom(roomID uint32) (*RoomManager, error) {
	roomsLock.Lock()
	defer roomsLock.Unlock()

	if noRoom {
		return nil, errors.New("no room exists")
	}

	if firstRoomID < nextRoomID {
		if firstRoomID > roomID || roomID >= nextRoomID {
			return nil, errors.New("no rooms exists")
		} 
	}

	if firstRoomID > nextRoomID {
		if nextRoomID <= roomID && roomID < firstRoomID {
			return nil, errors.New("no rooms exists")
		}
	}

	return roomManagers[roomID], nil
}

func AddRoom() (*RoomManager, error) {
	roomsLock.Lock()
	defer roomsLock.Unlock()

	if !noRoom && nextRoomID == firstRoomID {
		dur := time.Since(*roomManagers[firstRoomID].StartTime())
		if int64(dur) < maxGuranteedRoomDuration.Nanoseconds() {
			return nil, errors.New("all room spots filled")
		}
		roomManagers[firstRoomID].DeleteRoom()
		firstRoomID = (firstRoomID + 1) % maxRoomCount
	}

	newRoom := NewRoomManager(nextRoomID)
	roomManagers[nextRoomID] = newRoom

	nextRoomID = (nextRoomID + 1) % maxRoomCount
	noRoom = false

	return newRoom, nil
}
