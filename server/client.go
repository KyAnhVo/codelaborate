package main

import (
	"net"
)

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
