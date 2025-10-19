package main

import (
	"net"
	"sync"
)

// RoomManager class

var nextRoomID uint64 = 0

type RoomManager struct {
	roomID 			uint64
	latestVer		uint64
	clientCount 	uint8
	msgQueue 		*Queue[*UpdateMsg]
	client 			[]*Client
	lock			sync.RWMutex
}

func CreateRoom() *RoomManager {
	room 			:= new(RoomManager)
	room.lock.Lock()
	room.roomID 	= nextRoomID
	room.msgQueue 	= NewQueue[*UpdateMsg](0)
	room.client 	= make([]*Client, 255)
	nextRoomID++
	room.lock.Unlock()
	return room
}

// RoomID returns id of room (constant)
func (room *RoomManager) RoomID() uint64 {
	return room.roomID
}

// LatestVer returns latest version written by 
// the room, thread safe
func (room *RoomManager) LatestVer() uint64 {
	room.lock.RLock()
	defer room.lock.RUnlock()

	return room.latestVer
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

func (room *RoomManager) GetClient(index uint8) *Client {
	room.lock.RLock()
	defer room.lock.RUnlock()

	if index >= room.clientCount {
		return nil
	}
	return room.client[index]
}

func (room *RoomManager) AddClient(client *Client) bool {
	if room.clientCount == 255 {
		return false
	}
	room.client[room.clientCount] = client
	room.clientCount++
	return true
}

///////////////////////////////////////////////////////////////////////

// Utility types (message and client)

type CreateJoinMsg struct {
	Operation 		byte 	// 'C' or 'J'
	RoomID 			uint64
}


type UpdateMsg struct {
	ClientID		uint8	// partial key defined also by RoomID
	ClientVersion 	uint64	// version client is on, so we can know where to put it in prio Queue

	// delete [CursorPos, CursorPos + DeleteLen - 1]
	// then add InsertStr at CursorPos

	CursorPos		uint64
	DeleteLen		uint64
	InsertLen		uint64
	InsertStr		string
}

type Client struct {
	clientID 		uint8 	// partial key defined also by RoomID
	currentVer 		uint64
	connection 		net.Conn
}

func (client *Client) ID() uint8 {
	return client.clientID
}

func (client *Client) Connection() net.Conn {
	return client.connection
}

func (client *Client) setCurrentVer(version uint64) {
	client.currentVer = version
}

func (client *Client) getCurrentVer() uint64 {
	return client.currentVer
}
