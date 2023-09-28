package protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TESTING_RESPONSE_OK       = uint8(0x0)
	TESTING_RESPONSE_ERROR    = uint8(0xA)
	TESTING_TYPE_INT          = uint8(0x1)
	TESTING_TYPE_UNSIGNED_INT = uint8(0x2)
	TESTING_TYPE_STRING       = uint8(0x3)
)

func TestWriteError(t *testing.T) {
	response := NewErrorResponse("error_еррор")
	writer := WriterV1{}
	buf := writer.Write(response)
	reader := bytes.NewReader(buf)

	b, _ := reader.ReadByte()
	assert.Equal(t, TESTING_RESPONSE_ERROR, b)

	strBuf := buf[1:]
	assert.Equal(t, "error_еррор\r\n", string(strBuf))
}

func TestWriteSuccessWithIntValue(t *testing.T) {
	response := NewOkResponse(int64(1992))
	writer := WriterV1{}
	buf := writer.Write(response)
	reader := bytes.NewReader(buf)

	b, _ := reader.ReadByte() // Response type
	assert.Equal(t, TESTING_RESPONSE_OK, b)

	b, _ = reader.ReadByte() // Value type
	assert.Equal(t, TESTING_TYPE_INT, b)

	var val int64
	binary.Read(reader, binary.LittleEndian, &val)
	assert.Equal(t, int64(1992), val)
}

func TestWriteSuccessWithUnsignedIntValue(t *testing.T) {
	response := NewOkResponse(int64(-1992))
	writer := WriterV1{}
	buf := writer.Write(response)
	reader := bytes.NewReader(buf)

	b, _ := reader.ReadByte() // Response type
	assert.Equal(t, TESTING_RESPONSE_OK, b)

	b, _ = reader.ReadByte() // Value type
	assert.Equal(t, TESTING_TYPE_UNSIGNED_INT, b)

	var val int64
	binary.Read(reader, binary.LittleEndian, &val)
	assert.Equal(t, int64(-1992), val)
}

func TestWriteSuccessWithStringValue(t *testing.T) {
	response := NewOkResponse("hello_world")
	writer := WriterV1{}
	buf := writer.Write(response)
	reader := bytes.NewReader(buf)

	b, _ := reader.ReadByte() // Response type
	assert.Equal(t, TESTING_RESPONSE_OK, b)

	b, _ = reader.ReadByte() // Value type
	assert.Equal(t, TESTING_TYPE_STRING, b)

	strBuf := buf[2:]
	assert.Equal(t, "hello_world\r\n", string(strBuf))
}
