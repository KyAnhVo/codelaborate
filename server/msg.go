package main

type CreateJoinMsg struct {
	Operation 		byte 	// 'C' or 'J'
	RoomID 			uint32
}


type UpdateMsg struct {
	ClientID		uint8	// partial key defined also by RoomID
	// delete [CursorPos, CursorPos + DeleteLen - 1]
	// then add InsertStr at CursorPos
	closeconn		byte
	CursorPos		uint64
	DeleteLen		uint64
	InsertLen		uint64
	InsertStr		string
}

