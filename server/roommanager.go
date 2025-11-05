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

	if firstRoomID < nextRoomID {
		if !(firstRoomID <= roomID && roomID < nextRoomID) {
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
		dur := time.Now().Sub(*roomManagers[firstRoomID].StartTime())
		if int64(dur) < maxGuranteedRoomDuration.Nanoseconds() {
			return nil, errors.New("all room spots filled")
		}
		roomManagers[firstRoomID].DeleteRoom()
		firstRoomID = (firstRoomID + 1) % maxRoomCount
	}

	newRoom := NewRoomManager(nextRoomID)
	roomManagers[nextRoomID] = newRoom
	return newRoom, nil
}


///////////////////////////////////////////////////////////////////////

type RoomManager struct {
	roomID 			uint32
	clientCount 	uint8
	document		string
	msgQueue 		*Queue[*UpdateMsg]
	client 			[]*Client
	startTime 		time.Time
	lock			sync.RWMutex
	queueEdited		bool
}

// NewRoomManager creates a RoomManager for a new room.
func NewRoomManager(roomID uint32) *RoomManager {
	room := new(RoomManager)
	room.lock.Lock()
	room.roomID = roomID
	room.msgQueue = NewQueue[*UpdateMsg](0)
	room.document = ""
	room.client = make([]*Client, 255)
	room.lock.Unlock()
	room.startTime = time.Now()

	go room.RoomMainManager()
	return room
}

// Function to manage room
func (room *RoomManager) RoomMainManager() {
	for true {
		room.RoomMainManagerIteration()
	}
}

func (room *RoomManager) RoomMainManagerIteration() {
	room.lock.Lock()
	defer room.lock.Unlock()

	if !room.queueEdited {
		return
	}
	
	// look through the queue and then edit the queue
	for !room.msgQueue.IsEmpty() {
		msg, ok := room.msgQueue.Dequeue()
		if !ok {
			return
		}
		for _, client := range room.client {
			client.readChann <- msg
		}
	}
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
func (room *RoomManager) AddClient(conn net.Conn) (uint8, error) {
	if room.clientCount == 255 {
		return 0, errors.New("No available slot")
	}
	room.client[room.clientCount] = NewClient(room.clientCount, conn)
	room.clientCount++
	return room.clientCount - 1, nil
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
	readChann		chan *UpdateMsg
}

func NewClient(clientID uint8, c net.Conn) *Client {
	chann := make(chan *UpdateMsg)
	return &Client {
		clientID: clientID,
		connection: c,
		readChann: chann,
	}
}

func (client *Client) ID() uint8 {
	return client.clientID
}

func (client *Client) Connection() net.Conn {
	return client.connection
}

func (client *Client) ReadChan() chan *UpdateMsg {
	return client.readChann
}
