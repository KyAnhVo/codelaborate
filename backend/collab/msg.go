package collab

type JoinMsg struct {
	// 6 digits, sending user's session ID (not user ID)
	SessionID int `json:"SessionID"`

	// either "join" or "create"
	Op string `json:"Op"` 		

	// 6 digits, only matter for Op=="create"
	RoomID string `json:"RoomID"`	
}

type UpdateMsg struct {
	// 6 digits, sending user's session ID (not user ID)
	SessionID int `json:"SessionID"`

	// either "update" or "exit"
	Op string `json:"Op"`

	YjsBytes []byte
}
