package main

import (
	"net"
)

type Client struct {
	clientID 		uint8 	// partial key defined also by RoomID
	connection 		net.Conn
	roomManager 	*RoomManager
	writeQueue		*Queue[*UpdateMsg]
	readChann		chan *UpdateMsg
}

func NewClient(clientID uint8, c net.Conn, room *RoomManager, writeQueue *Queue[*UpdateMsg]) *Client {
	chann := make(chan *UpdateMsg)
	return &Client {
		clientID: clientID,
		connection: c,
		readChann: chann,
		writeQueue: writeQueue,
		roomManager: room,
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

func (client *Client) RoomManager() *RoomManager {
	return client.roomManager
}
