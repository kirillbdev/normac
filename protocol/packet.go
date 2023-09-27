package protocol

const (
	CMD_EMPTY = 0
	CMD_PING  = 1
	CMD_GET   = 2
	CMD_SET   = 3
	CMD_DEBUG = 255
)

type Command uint8

type Packet struct {
	Command     Command
	Key         string
	Value       any
	PingMessage string
}

func NewEmptyV1Packet() *Packet {
	return &Packet{
		Command: CMD_EMPTY,
		Key:     "",
		Value:   "",
	}
}
