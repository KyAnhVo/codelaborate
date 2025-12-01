package main

import (
	"errors"
	"net"
	"sync"
	"io"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
)

// HandleConnection handles a connection from the client and 
// processes that connection
// If there is an error, send code 3 to client. Else send code 2.
func HandleConnection(wg *sync.WaitGroup, c net.Conn) {
	errorCode := make([]byte, 1)
	CodeOk := byte(2)
	CodeEr := byte(3)

	errorCode[0] = CodeEr

	defer wg.Done()
	joinMsg := GetConnection(c)

	room, err := ProcessRoomRequest(joinMsg)
	if err != nil {
		c.Write(errorCode)
		c.Close()
		return
	}

	cliID, err := room.AddClient(c)
	if err != nil {
		c.Write(errorCode)
		c.Close()
		return
	}
	client := room.GetClient(cliID)
	log.Infof("Client:\tID: %d\tROOM: %d", client.clientID, client.roomManager.roomID)
	
	// send error code CodeOk + roomID
	errorCode[0] = CodeOk
	c.Write(errorCode)
	err = sendUint32(c, room.roomID)
	if err != nil {
		log.Errorf("Failed to send roomID: %v", err)
		c.Close()
		return
	}
	err = sendUint8(c, client.clientID)
	if err != nil {
		log.Errorf("Failed to send clientID: %v", err)
		c.Close()
		return
	}

	go ConnToRoomManager(client)
	RoomManagerToConn(client)
}

// GetConnection receives join or create operation msg.
// Binary msg layout (big-endian):
// 	[0-0]  	uint8_t 	operation 	- Operation, either 'C' for create or 'J' for join
// 	[1-4] 	uint32_t 	roomId 		- Id of room, ignored if operation == 'C' 
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
// 	[0-0]		1 if close connection, 0 if send text
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
		var err error
		msg := new(UpdateMsg)
		_, err = io.ReadFull(conn, byteBuffer)
		if err != nil {
			client.connection.Close()
			return
		}
		switch byteBuffer[0] {
			case CLOSECONN:
				msg.closeconn = CLOSECONN
				continue
			case UPDATE:
				msg.closeconn = UPDATE
			default:
				continue
		}

		_, err = io.ReadFull(conn, uint64Buffer)
		if err != nil {
			client.connection.Close()
			return
		}
		msg.CursorPos = binary.BigEndian.Uint64(uint64Buffer)

		_, err = io.ReadFull(conn, uint64Buffer)
		if err != nil {
			client.connection.Close()
			return
		}
		msg.DeleteLen = binary.BigEndian.Uint64(uint64Buffer)

		_, err = io.ReadFull(conn, uint64Buffer)
		if err != nil {
			client.connection.Close()
			return
		}
		msg.InsertLen = binary.BigEndian.Uint64(uint64Buffer)

		strBuffer = make([]byte, msg.InsertLen)
		_, err = io.ReadFull(conn, strBuffer)
		if err != nil {
			client.connection.Close()
			return
		}
		msg.InsertStr = string(strBuffer)

		msg.ClientID = client.clientID

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

		var err error
		// client.connection.Write(closeconn)
		err = writeAll(client.connection, closeconn)
		if err != nil {
			log.Errorf("send failed: closeconn")
		}
		// client.connection.Write(cursorPos)
		err = writeAll(client.connection, cursorPos)
		if err != nil {
			log.Errorf("send failed: closeconn")
		}
		// client.connection.Write(deleteLen)
		err = writeAll(client.connection, deleteLen)
		if err != nil {
			log.Errorf("send failed: closeconn")
		}
		// client.connection.Write(insertLen)
		err = writeAll(client.connection, insertLen)
		if err != nil {
			log.Errorf("send failed: closeconn")
		}
		// client.connection.Write(insertStr)
		err = writeAll(client.connection, insertStr)
		if err != nil {
			log.Errorf("send failed: closeconn")
		}
	}
}

func writeAll(conn net.Conn, data []byte) error {
    written := 0
    for written < len(data) {
        n, err := conn.Write(data[written:])
        if err != nil {
            return err
        }
        written += n
    }
    return nil
}

func sendUint32(c net.Conn, num uint32) error {
	bytearr := make([]byte, 4)
	binary.BigEndian.PutUint32(bytearr, num)
	return writeAll(c, bytearr)
}

func sendUint8(c net.Conn, num uint8) error {
	bytearr := make([]byte, 1)
	bytearr[0] = num
	return writeAll(c, bytearr)
}

// -------------------------------------------------------------------------


