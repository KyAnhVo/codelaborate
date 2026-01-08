package collab

type JoinMsg struct {
	// 6 digits, sending user's session ID (not user ID)
	SessionID string `json:"SessionID"`

	// either "join" or "create"
	Op string `json:"Op"` 		

	// 6 digits, only matter for Op=="create"
	RoomID string `json:"RoomID"`	
}

type UpdateMsg struct {
	// 6 digits, sending user's session ID (not user ID)
	SessionID string `json:"SessionID"`

	// either "update" or "exit"
	Op string `json:"Op"`

	// the entire below matters when Op=="update"
	// the system contains CursorPos, DeleteLen,
	// InsertLen, and InsertStr.
	// When delete, 	InsertLen == 0 and DeleteLen != 0
	// When insert, 	InsertLen != 0 and DeleteLen == 0
	// when replace,	InsertLen != 0 and DeleteLen != 0

	// Length of deleted string
	DeleteLen int `json:"DeleteLen"`

	// Length of added string
	InsertLen int `json:"InsertLen"`

	// Length of inserted bytes
	InsertStr []byte `json:"InsertStr"`
}
