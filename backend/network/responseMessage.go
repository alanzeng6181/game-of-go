package network

type ResponseMessage struct {
	ResponseType ResponseType
	Content      string
}

type ResponseType string

const (
	State   ResponseType = "State"
	Message ResponseType = "Message"
	Error   ResponseType = "Error"
)
