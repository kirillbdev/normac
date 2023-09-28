package protocol

import (
	"encoding/binary"
)

type WriterV1 struct {
}

func (w *WriterV1) Write(response *Response) []byte {
	var buf []byte

	switch response.ResponseType {
	case RESPONSE_OK:
		buf = w.writeByte(buf, 0x0)
	case RESPONSE_ERR:
		buf = w.writeByte(buf, 0xA)
	}

	if response.ResponseType == RESPONSE_ERR {
		buf = w.writeString(buf, response.ErrorMessage)
	} else {
		switch response.Value.(type) {
		case int64:
			if response.Value.(int64) >= 0 {
				buf = w.writeByte(buf, 0x1)
			} else {
				buf = w.writeByte(buf, 0x2)
			}
			buf = w.writeInt(buf, response.Value.(int64))
		case string:
			buf = w.writeByte(buf, 0x3)
			buf = w.writeString(buf, response.Value.(string))
		}
	}

	return buf
}

func (w *WriterV1) writeByte(buf []byte, val byte) []byte {
	return append(buf, val)
}

func (w *WriterV1) writeInt(buf []byte, val int64) []byte {
	return binary.LittleEndian.AppendUint64(buf, uint64(val))
}

func (w *WriterV1) writeString(buf []byte, val string) []byte {
	strBuf := []byte(val + "\r\n")

	return append(buf, strBuf...)
}
