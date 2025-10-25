package main

import (
	"net"
	"sync"
	"io"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
	"errors"
)

func HandleConnection(wg *sync.WaitGroup, c net.Conn) {
	defer wg.Done()
	var room 	*RoomManager
	var err 	error

	joinMsg := GetConnectionMsg(c)


	
}

// GetConnectionMsg receives join or create operation msg.
// Binary msg layout (big-endian):
// 	[0-5]  	uint8_t 	operation 	- Operation, either 'C' for create or 'J' for join
// 	[1-5] 	uint32_t 	roomId 		- Id of room, ignored if operation == 'C' 
// Message length = 5:
func GetConnectionMsg(c net.Conn) *CreateJoinMsg {
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

func GetConnection(joinMsg *CreateJoinMsg) (*RoomManager, error) {
	if joinMsg.Operation == 'C' {
		log.Info("Msg: Create room")
	} else {
		log.Infof("Msg: Join room, id: %d", joinMsg.RoomID)
	}

	if joinMsg.Operation == 'C' {
		room, err := AddRoom()
		if err != nil {
			log.Errorf("Cannot add any more rooms")
			return nil, errors.New("Can't add room")
		}
		return room, nil
	} else {
		room, err := JoinRoom(joinMsg.RoomID)
		if err != nil {
			log.Errorf("RoomID invalid")
			return nil, errors.New("No room with such ID")
		}
		return room, nil
	}
}
