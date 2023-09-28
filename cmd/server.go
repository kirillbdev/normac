package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/kirillbdev/normac/protocol"
	"github.com/kirillbdev/normac/storage"
	"log"
	"net"
	"strconv"
)

const (
	PROTOCOL_VERSION_1 = 0x1
)

type Session struct {
	conn   net.Conn
	writer protocol.Writer
}

func (s *Session) Read() (*bytes.Reader, error) {
	var buffer = make([]byte, 1024*2) // Packet size = 2MB, hardcoded by default yet

	_, err := s.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buffer), nil
}

func (s *Session) Send(response *protocol.Response) {
	s.conn.Write(s.writer.Write(response))
	s.conn.Close()
}

func getProtoReader(r *bytes.Reader) (protocol.Reader, error) {
	v, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch v {
	case PROTOCOL_VERSION_1:
		return protocol.NewReaderV1(), nil
	default:
		return nil, errors.New("Undefined protocol version")
	}
}

func handleRequest(session *Session, storage *storage.Storage) {
	reader, err := session.Read()
	if err != nil {
		session.Send(protocol.NewErrorResponse("[CONN_ERR] " + err.Error()))
		return
	}

	protoReader, err := getProtoReader(reader)
	if err != nil {
		session.Send(protocol.NewErrorResponse("[REQ_ERR] " + err.Error()))
		return
	}

	var response *protocol.Response

	packet, err := protoReader.Read(reader)
	if err != nil {
		session.Send(protocol.NewErrorResponse(err.Error()))
		return
	}

	switch packet.Command {
	case protocol.CMD_PING:
		response = protocol.NewOkResponse("PONG")
		if len(packet.PingMessage) > 0 {
			response.Value = response.Value.(string) + " " + packet.PingMessage
		}
	case protocol.CMD_GET:
		val, ok := storage.Get(packet.Key)
		if ok {
			response = protocol.NewOkResponse(val)
		} else {
			response = protocol.NewErrorResponse("Key not found")
		}
	case protocol.CMD_SET:
		storage.Set(packet.Key, packet.Value)
		response = protocol.NewOkResponse(packet.Value)
	case protocol.CMD_DEBUG:
		fmt.Println(storage.Data())
		response = protocol.NewOkResponse("OK")
	default:
		response = protocol.NewErrorResponse("Invalid command")
	}

	session.Send(response)
}

func Run(port int) {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	// Hardcoded initialize yet
	strg := storage.NewStorage(1000)
	w := protocol.WriterV1{}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}

		session := Session{
			conn:   conn,
			writer: &w,
		}

		go handleRequest(&session, strg)
	}
}
