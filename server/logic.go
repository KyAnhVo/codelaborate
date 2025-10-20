package main

import (
	"net"
	"sync"
	"io"
	"encoding/binary"
	"fmt"
)

func HandleConnection(wg *sync.WaitGroup, c net.Conn) {
	defer wg.Done()
	
	joinMsg := GetConnection(c)
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
