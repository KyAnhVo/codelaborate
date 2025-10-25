package main

import (
	"net"
	"sync"
	"errors"
	"time"
)

// Only 128 rooms are available at a time max.
const maxRoomCount uint32 = 128

// Duration room is guranteed to exist before being available to be kicked out
const maxGuranteedRoomDuration = 30 * time.Minute

///////////////////////////////////////////////////////////////////////

// Each room will get its ID incrementally from here.

var roomManagers []*RoomManager
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

	if firstRoomID == nextRoomID {
		return roomManagers[roomID], nil
	}

	// no wraparound, check sides
	if firstRoomID < nextRoomID {
		if firstRoomID <= roomID && roomID < nextRoomID {
			return roomManagers[roomID], nil
		} 
		return nil, errors.New("no rooms exists")
	}

	// wraparound, check middle
	if nextRoomID <= roomID && roomID < firstRoomID {
		return nil, errors.New("no rooms exists")
	}
	return roomManagers[roomID], nil
}

// AddRoom attempts to create a room, but rejects (returns error)
// if max amount of room is exceeded and non of the room
// has exceeded their guaranteed time.
func AddRoom() (*RoomManager, error) {
	roomsLock.Lock()
	defer roomsLock.Unlock()

	// if there is no room left available, we will check the first
	// created room. 
	// If such room exists longer than the guaranteed time,
	// we terminate that room.
	// Else we deny the room creation function.
	if !noRoom && nextRoomID == firstRoomID {
		dur := time.Since(*roomManagers[firstRoomID].StartTime())
		if int64(dur) < maxGuranteedRoomDuration.Nanoseconds() {
			return nil, errors.New("all room spots filled")
		}
		roomManagers[firstRoomID].DeleteRoom()
		firstRoomID = (firstRoomID + 1) % maxRoomCount
	}
	
	// Create room, increment nextRoomID
	newRoom := NewRoomManager(nextRoomID)
	roomManagers[nextRoomID] = newRoom
	nextRoomID = (nextRoomID + 1) % maxRoomCount
	return newRoom, nil
}

// TODO: implement CreateRoom function

///////////////////////////////////////////////////////////////////////

type RoomManager struct {
	roomID 			uint32
	clientCount 	uint8
	document		string
	msgQueue 		*Queue[*UpdateMsg]
	client 			[]*Client
	startTime 		time.Time
	lock			sync.RWMutex
}

// NewRoomManager creates a RoomManager for a new room.
func NewRoomManager(roomID uint32) *RoomManager {
	room 			:= new(RoomManager)
	room.lock.Lock()
	room.roomID 	= roomID
	room.msgQueue 	= NewQueue[*UpdateMsg](0)
	room.document 	= ""
	room.client 	= make([]*Client, 255)
	room.lock.Unlock()
	room.startTime 		= time.Now()
	return room
}

// RoomID returns id of room (constant)
func (room *RoomManager) RoomID() uint32 {
	return room.roomID
}

// Document returns the stored document for this room
func (room *RoomManager) Document() string {
	room.lock.RLock()
	defer room.lock.RUnlock()

	return room.document
}

// EnqueueMsg enqueues an editing message waiting for 
// processing.
func (room *RoomManager) EnqueueMsg(msg *UpdateMsg) {
	room.lock.Lock()
	defer room.lock.Unlock()

	room.msgQueue.Enqueue(msg)
}

// DequeueMsg dequeues a message for processing, or nil if
// queue is empty.
func (room *RoomManager) DequeueMsg() *UpdateMsg {
	room.lock.Lock()
	defer room.lock.Unlock()

	msg, ok := room.msgQueue.Dequeue()
	if ok {
		return msg
	}
	return nil
}

// GetClient gets the client given the client ID.
func (room *RoomManager) GetClient(ID uint8) *Client {
	room.lock.RLock()
	defer room.lock.RUnlock()

	if ID >= room.clientCount {
		return nil
	}
	return room.client[ID]
}

// AddClient adds the client that uses such connection
func (room *RoomManager) AddClient(conn net.Conn) bool {
	if room.clientCount == 255 {
		return false
	}
	room.client[room.clientCount] = NewClient(room.clientCount, conn)
	room.clientCount++
	return true
}

func (room *RoomManager) StartTime() *time.Time {
	return &room.startTime
}

// DeleteRoom closes all connections, somehow signal all threads 
// to remove its conn and kill themselves.
func (room *RoomManager) DeleteRoom() {
	// TODO: Implement DeleteRoom
}

///////////////////////////////////////////////////////////////////////

// Message types (pure data types)

type CreateJoinMsg struct {
	Operation 		byte 	// 'C' or 'J'
	RoomID 			uint32
}


type UpdateMsg struct {
	ClientID		uint8	// partial key defined also by RoomID
	// delete [CursorPos, CursorPos + DeleteLen - 1]
	// then add InsertStr at CursorPos
	CursorPos		uint64
	DeleteLen		uint64
	InsertLen		uint64
	InsertStr		string
}

///////////////////////////////////////////////////////////////////////

// Client type

type Client struct {
	clientID 		uint8 	// partial key defined also by RoomID
	connection 		net.Conn
}

func NewClient(clientID uint8, c net.Conn) *Client {
	return &Client {
		clientID: clientID,
		connection: c,
	}
}

func (client *Client) ID() uint8 {
	return client.clientID
}

func (client *Client) Connection() net.Conn {
	return client.connection
}

