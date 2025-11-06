package main

import (
	"net"
)

type Client struct {
	clientID 		uint8 	// partial key defined also by RoomID
	connection 		net.Conn
	writeQueue		*Queue[*UpdateMsg]
	readChann		chan *UpdateMsg
}

func NewClient(clientID uint8, c net.Conn, writeQueue *Queue[*UpdateMsg]) *Client {
	chann := make(chan *UpdateMsg)
	return &Client {
		clientID: clientID,
		connection: c,
		readChann: chann,
		writeQueue: writeQueue,
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
