package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
)

type ReaderV1 struct {
}

func NewReaderV1() *ReaderV1 {
	return &ReaderV1{}
}

func (decoder *ReaderV1) Read(reader *bytes.Reader) (*Packet, error) {
	packet := NewEmptyV1Packet()

	// Step 1: read first byte - command
	b, err := reader.ReadByte()
	if err != nil {
		return nil, errors.New("Reader error")
	}

	cmd := int(b)
	switch cmd {
	case CMD_PING, CMD_SET, CMD_GET:
		packet.Command = Command(cmd)
	default:
		return nil, errors.New("Undefined command " + strconv.Itoa(cmd))
	}

	if cmd == CMD_PING {
		str, ok := decoder.readString(reader)
		if ok {
			packet.PingMessage = str
		}
	} else if cmd == CMD_GET {
		key, ok := decoder.readString(reader)
		if !ok {
			return nil, errors.New("Expected key (string)")
		}

		packet.Key = key
	} else if cmd == CMD_SET {
		key, ok := decoder.readString(reader)
		if !ok {
			return nil, errors.New("Expected key (string)")
		}
		packet.Key = key

		// Read type
		tp, err := reader.ReadByte()
		if err != nil {
			return nil, errors.New("Expected value type (byte)")
		}

		switch int(tp) {
		case 1: // INT
			val, ok := decoder.readInt(reader)
			if !ok {
				return nil, errors.New("Expected value type (int)")
			}

			packet.Value = val
		case 3: // STRING
			val, ok := decoder.readString(reader)
			if !ok {
				return nil, errors.New("Expected value type (string)")
			}

			packet.Value = val
		default:
			return nil, errors.New("Incorect value type")
		}

	}

	return packet, nil
}

func (decoder *ReaderV1) readInt(reader *bytes.Reader) (int64, bool) {
	var res int64

	err := binary.Read(reader, binary.LittleEndian, &res)
	if err != nil {
		return 0, false
	}

	return res, true
}

func (decoder *ReaderV1) readString(reader *bytes.Reader) (string, bool) {
	result := ""
	foundEnd := false

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "", false
		}

		if b == 0x0 {
			break
		} else if b == 0x0D { // \r
			// Try to expect \n in next byte
			next, err := reader.ReadByte()
			if err != nil {
				return "", false
			}

			if next == 0x0A {
				foundEnd = true
				break
			}
		} else {
			result += string(b)
		}
	}

	if foundEnd && len(result) > 0 {
		return result, true
	}

	return "", false
}
