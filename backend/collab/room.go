package collab

type Room struct {
	// ID of room (multiple rooms possible)
	Id string

	// all msgs sent into room goes through this
	MsgInChannel chan *UpdateMsg

	// this manages client, send update to clients, etc.
	Clients []*Client
}

func CreateRoom(roomId string) *Room {
	room := new(Room)
	room.Id = roomId
	room.MsgInChannel = make(chan *UpdateMsg, 128)
	room.Clients = make([]*Client, 32)
	return room
}

func (r *Room) Start() {
	for {
		
	}
}

func (r *Room) processMsg() {

}
