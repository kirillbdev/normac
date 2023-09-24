package protocol

import "bytes"

type Reader interface {
	Read(reader *bytes.Reader) (*Packet, error)
}
