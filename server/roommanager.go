package main

import (
	"net"
	"sync"
	"errors"
	"time"
)

type RoomManager struct {
	roomID 			uint32
	clientCount 	uint8
	document		string
	msgQueue 		*Queue[*UpdateMsg]
	client 			[]*Client
	startTime 		time.Time
	queueEdited		bool

	lock			sync.RWMutex
	lockSignal		*sync.Cond
}

// NewRoomManager creates a RoomManager for a new room.
func NewRoomManager(roomID uint32) *RoomManager {
	room := new(RoomManager)
	room.lockSignal = sync.NewCond(&room.lock)

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

// RoomMainManager is supposed to be threaded
// out, gets all the update messages going in,
// change them respectively and sends back to
// all the client threads.
func (room *RoomManager) RoomMainManager() {
	room.lock.Lock()
	defer room.lock.Unlock()

	for true {
		// sleep, woken up when a msg is enqueued.
		room.lockSignal.Wait()

		// look through the queue and then edit the msg and send
		// to all the threads
		for !room.msgQueue.IsEmpty() {
			msg, ok := room.msgQueue.Dequeue()
			if !ok {
				return
			}
			for _, client := range room.client {
				if client == nil {
					continue
				}

				if msg.ClientID == client.clientID {
					continue
				}
				client.readChann <- msg
			}
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
	room.lockSignal.Signal()
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
	room.client[room.clientCount] = NewClient(room.clientCount, conn, room.msgQueue)
	room.clientCount++
	return room.clientCount - 1, nil
}

// StartTime returns the time since the room was started
func (room *RoomManager) StartTime() *time.Time {
	return &room.startTime
}

// DeleteRoom closes all connections, somehow signal all threads 
// to remove its conn and kill themselves.
func (room *RoomManager) DeleteRoom() {
	// TODO: Implement DeleteRoom
	// Bruh where TF am I even starting at?
}
