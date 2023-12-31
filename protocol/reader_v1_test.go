package protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	BYTE_CMD_SET           = 0x3
	TEST_TYPE_POSITIVE_INT = 1
	TEST_TYPE_NEGATIVE_INT = 2
)

func TestReadPingCommand(t *testing.T) {
	buf := []byte{0x1, 0x1}
	reader := NewReaderV1()
	packet, _ := reader.Read(bytes.NewReader(buf))

	assert.Equal(t, Command(CMD_PING), packet.Command)
}

func TestReadPingWithMessageCommand(t *testing.T) {
	msg := []byte("Hello world\r\n")
	buf := append([]byte{0x1}, msg...)

	reader := NewReaderV1()
	packet, _ := reader.Read(bytes.NewReader(buf))

	assert.Equal(t, Command(CMD_PING), packet.Command)
	assert.Equal(t, "Hello world", packet.PingMessage)
}

func TestGetCommand(t *testing.T) {
	msg := []byte("some_key\r\n")
	buf := append([]byte{0x2}, msg...)

	reader := NewReaderV1()
	packet, _ := reader.Read(bytes.NewReader(buf))

	assert.Equal(t, Command(CMD_GET), packet.Command)
	assert.Equal(t, "some_key", packet.Key)
}

func TestSetCommandWithUnsignedIntValue(t *testing.T) {
	key := []byte("some_key\r\n")
	val := make([]byte, 8)
	binary.LittleEndian.PutUint64(val, 1992)
	buf := append([]byte{BYTE_CMD_SET}, key...)
	buf = append(buf, []byte{TEST_TYPE_POSITIVE_INT}...)
	buf = append(buf, val...)

	reader := NewReaderV1()
	packet, _ := reader.Read(bytes.NewReader(buf))

	assert.Equal(t, Command(CMD_SET), packet.Command)
	assert.Equal(t, "some_key", packet.Key)
	assert.IsType(t, int64(1), packet.Value)
	assert.Equal(t, int64(1992), packet.Value)
}

func TestSetCommandWithSignedIntValue(t *testing.T) {
	key := []byte("some_key\r\n")
	val := make([]byte, 8)
	binary.LittleEndian.PutUint64(val, 1992)
	buf := append([]byte{BYTE_CMD_SET}, key...)
	buf = append(buf, []byte{TEST_TYPE_NEGATIVE_INT}...)
	buf = append(buf, val...)

	reader := NewReaderV1()
	packet, _ := reader.Read(bytes.NewReader(buf))

	assert.Equal(t, Command(CMD_SET), packet.Command)
	assert.Equal(t, "some_key", packet.Key)
	assert.IsType(t, int64(1), packet.Value)
	assert.Equal(t, int64(-1992), packet.Value)
}

func TestSetCommandWithStringValue(t *testing.T) {
	key := []byte("some_key\r\n")
	val := []byte("hello world\r\n")
	buf := append([]byte{BYTE_CMD_SET}, key...)
	buf = append(buf, 0x3)
	buf = append(buf, val...)

	reader := NewReaderV1()
	packet, _ := reader.Read(bytes.NewReader(buf))

	assert.Equal(t, Command(CMD_SET), packet.Command)
	assert.Equal(t, "some_key", packet.Key)
	assert.IsType(t, "", packet.Value)
	assert.Equal(t, "hello world", packet.Value)
}

func TestSetCommandWithUtf8StringValue(t *testing.T) {
	key := []byte("some_key\r\n")
	val := []byte("привет мир\r\n")
	buf := append([]byte{BYTE_CMD_SET}, key...)
	buf = append(buf, 0x3)
	buf = append(buf, val...)

	reader := NewReaderV1()
	packet, _ := reader.Read(bytes.NewReader(buf))

	assert.Equal(t, Command(CMD_SET), packet.Command)
	assert.Equal(t, "some_key", packet.Key)
	assert.IsType(t, "", packet.Value)
	assert.Equal(t, "привет мир", packet.Value)
}
