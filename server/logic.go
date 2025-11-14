package main

import (
	"errors"
	"net"
	"sync"
	"io"
	"encoding/binary"
)

// HandleConnection handles a connection from the client and 
// processes that connection
func HandleConnection(wg *sync.WaitGroup, c net.Conn) {
	defer wg.Done()
	joinMsg := GetConnection(c)

	room, err := ProcessRoomRequest(joinMsg)
	if err != nil {
		io.WriteString(c, err.Error())
	}

	cliID, err := room.AddClient(c)
	if err != nil {
		io.WriteString(c, err.Error())
	}
	client := room.GetClient(cliID)

	go ConnToRoomManager(client)
	RoomManagerToConn(client)
}

// GetConnection receives join or create operation msg.
// Binary msg layout (big-endian):
// 	[0-5]  	uint8_t 	operation 	- Operation, either 'C' for create or 'J' for join
// 	[1-5] 	uint32_t 	roomId 		- Id of room, ignored if operation == 'C' 
// Message length = 5:
func GetConnection(c net.Conn) *CreateJoinMsg {
	opBuffer := make([]byte, 1)
	idBuffer := make([]byte, 4)

	_, err := io.ReadFull(c, opBuffer)
	if err != nil {
		return nil
	}
	_, err = io.ReadFull(c, idBuffer)
	if err != nil {
		return nil
	}

	roomID := binary.BigEndian.Uint32(idBuffer)

	return &CreateJoinMsg{
		Operation: 	opBuffer[0],
		RoomID: 	roomID,	
	}
}

func ProcessRoomRequest(msg *CreateJoinMsg) (*RoomManager, error) {
	switch msg.Operation {
	case 'C':
		return AddRoom()
	case 'J':
		return JoinRoom(msg.RoomID)
	default:
		return nil, errors.New("invalid operation")
	}
}

// ConnToRoomManager gets the following text from client:
// 	[0-0]		0 if close connection, 1 if send text
//	[1-8]		cursor position (0-indexed)
//	[9-16]		delete length
//	[17-24]		insert length
//	[25-...]	insert string (not null terminated)
func ConnToRoomManager(client *Client) {
	CLOSECONN := byte(1)
	UPDATE := byte(0)
	conn := client.Connection()
	byteBuffer := make([]byte, 1)
	uint64Buffer := make([]byte, 8)
	var strBuffer []byte
	for {
		msg := new(UpdateMsg)
		io.ReadFull(conn, byteBuffer)
		switch byteBuffer[0] {
		case CLOSECONN:
			msg.closeconn = CLOSECONN
			continue
		case UPDATE:
			msg.closeconn = UPDATE
		default:
			continue
		}

		io.ReadFull(conn, uint64Buffer)
		msg.CursorPos = binary.BigEndian.Uint64(uint64Buffer)
		io.ReadFull(conn, uint64Buffer)
		msg.DeleteLen = binary.BigEndian.Uint64(uint64Buffer)
		io.ReadFull(conn, uint64Buffer)
		msg.InsertLen = binary.BigEndian.Uint64(uint64Buffer)

		strBuffer = make([]byte, msg.InsertLen)
		io.ReadFull(conn, strBuffer)
		msg.InsertStr = string(strBuffer)

		client.roomManager.EnqueueMsg(msg)
	}
}

// RoomManagerToConn sends msgs to clients
//	[0-0]		0 if close connection, 1 if update message
// 	[1-8]		cursor position for edit
//	[9-16]		delete length from cursor
//	[17-24]		length of string to insert into cursor pos
//	[25-...]	string to insert (not null terminated)
func RoomManagerToConn(client *Client) {
	closeconn := make([]byte, 1)
	cursorPos := make([]byte, 8)
	deleteLen := make([]byte, 8)
	insertLen := make([]byte, 8)
	var insertStr []byte
	for {
		msg := <- client.readChann
		closeconn[0] = msg.closeconn
		binary.BigEndian.PutUint64(cursorPos, msg.CursorPos)
		binary.BigEndian.PutUint64(deleteLen, msg.DeleteLen)
		binary.BigEndian.PutUint64(insertLen, msg.InsertLen)
		insertStr = []byte(msg.InsertStr)

		client.connection.Write(closeconn)
		client.connection.Write(cursorPos)
		client.connection.Write(deleteLen)
		client.connection.Write(insertLen)
		client.connection.Write(insertStr)
	}
}

// -------------------------------------------------------------------------


