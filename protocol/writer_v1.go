package protocol

import "encoding/binary"

type WriterV1 struct {
}

func (w *WriterV1) Write(response *Response) []byte {
	buf := make([]byte, 1024)
	buf = w.writeInt(buf, int64(response.ResponseType))

	if response.ResponseType == RESPONSE_ERR {
		buf = w.writeString(buf, response.ErrorMessage)
	} else {
		switch response.Value.(type) {
		case int64:
			buf = w.writeInt(buf, response.Value.(int64))
		case string:
			buf = w.writeString(buf, response.Value.(string))
		}
	}

	return buf
}

func (w *WriterV1) writeInt(buf []byte, val int64) []byte {
	intBUf := make([]byte, 8)
	binary.LittleEndian.PutUint64(intBUf, uint64(val))

	return append(buf, intBUf...)
}

func (w *WriterV1) writeString(buf []byte, val string) []byte {
	strBuf := []byte(val + "\r\n")

	return append(buf, strBuf...)
}
