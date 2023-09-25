package cmd

import (
	"bytes"
	"fmt"
	"github.com/kirillbdev/normac/protocol"
	"log"
	"net"
	"strconv"
	"sync"
)

type Storage struct {
	sync.RWMutex
	data map[string]any
}

func (storage *Storage) Set(key string, value any) {
	storage.Lock()
	defer storage.Unlock()
	storage.data[key] = value
}

func (storage *Storage) Get(key string) (any, bool) {
	storage.RLock()
	defer storage.RUnlock()

	val, ok := storage.data[key]

	return val, ok
}

var storage = Storage{
	data: make(map[string]any, 1024),
}

func handleRequest(conn net.Conn) {
	var buffer = make([]byte, 1024*16) // Packet size = 16MB, hardcoded by default yet
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// Step 1: read first byte - protocol version
	reader := bytes.NewReader(buffer)
	v, err := reader.ReadByte()
	if err != nil {
		log.Fatal(err)
	}

	if v != 0x1 {
		conn.Write([]byte("Invalid protocol version" + "\r\n"))
		conn.Close()
		return
	}

	var response protocol.Response

	protoReader := protocol.NewReaderV1()
	packet, err := protoReader.Read(reader)
	if err != nil {
		response = protocol.Response{
			ResponseType: protocol.RESPONSE_ERR,
			ErrorMessage: err.Error(),
		}
	} else {
		switch packet.Command {
		case protocol.CMD_PING:
			response = protocol.Response{
				ResponseType: protocol.RESPONSE_OK,
				Value:        "PONG",
			}

			if len(packet.Key) > 0 {
				response.Value = response.Value.(string) + " " + packet.Key
			}
		case protocol.CMD_GET:
			val, ok := storage.Get(packet.Key)
			if ok {
				response = protocol.Response{
					ResponseType: protocol.RESPONSE_OK,
					Value:        val,
				}
			} else {
				response = protocol.Response{
					ResponseType: protocol.RESPONSE_ERR,
					ErrorMessage: "Key not found",
				}
			}
		case protocol.CMD_SET:
			storage.Set(packet.Key, packet.Value)
			response = protocol.Response{
				ResponseType: protocol.RESPONSE_OK,
				Value:        packet.Value,
			}
		default:
			response = protocol.Response{
				ResponseType: protocol.RESPONSE_ERR,
				ErrorMessage: fmt.Sprintf("Invalid command %d", packet.Command),
			}
		}
	}

	w := protocol.WriterV1{}
	conn.Write(w.Write(&response))
	conn.Close()
}

func Run(port int) {
	listen, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleRequest(conn)
	}
}
