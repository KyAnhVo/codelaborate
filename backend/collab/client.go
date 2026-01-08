package collab

type Client struct {
	SessionID string
	OutChannel chan *UpdateMsg 
}
