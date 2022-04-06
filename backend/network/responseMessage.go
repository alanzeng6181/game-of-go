package network

type ResponseMessage struct {
	ResponseType ResponseType
	Content      string
}

type ResponseType int8

const (
	State ResponseType = iota
	Message
)
