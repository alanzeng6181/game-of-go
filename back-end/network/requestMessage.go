package network

type RequestMessage struct {
	Command   Command
	Arguments []string
}

type Command int8

const (
	Move Command = iota
	TakeBack
	Resign
	ClaimTerritory
	MessageCommand
)
