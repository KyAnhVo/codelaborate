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

	// variable declarations (You know, C style)
	var uint32ByteBuffer 	[]byte
	var uint8ByteBuffer 	[]byte
	var msgBuffer 			[]byte
	var roomID 				uint32

	// tmp buffer memory allocation
	uint32ByteBuffer 	= make([]byte, 4)
	uint8ByteBuffer 	= make([]byte, 1)
	msgBuffer 			= make([]byte, 1024)

	// Dev time: for unused vars (bruh i swear go gotta not do this)
	_ = msgBuffer

	// Handles Join or Create operation
	// 
	// Binary msg layout (big-endian):
	//
	// Field layout (total 5 bytes):
	// [0-0] 	uint8_t 	operation 	- value can be 'C' for Create or 'J' for Join
	// [1-5] 	uint32_t 	roomId 		- Id of room, ignored if operation == 'C' 
	//
	// Message length = 5:
	
	// read op
	_, err := io.ReadFull(c, uint8ByteBuffer)
	if err != nil {
		// TODO: implement error handling
		return
	}

	// read roomId
	_, err = io.ReadFull(c, uint32ByteBuffer)
	if err != nil {
		// TODO: implement error handling
		return
	}
	roomID = binary.BigEndian.Uint32(uint32ByteBuffer)
	fmt.Printf("%s: %d", uint8ByteBuffer, roomID)

	// process room_id


}
