package cmd

import (
	"bytes"
	"fmt"
	"github.com/kirillbdev/normac/protocol"
	"log"
	"net"
	"strconv"
)

var storage map[string]any = make(map[string]any, 1024)

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

	var response string

	protoReader := protocol.NewV1Reader()
	packet, err := protoReader.Read(reader)
	if err != nil {
		response = "[ERROR] " + err.Error()
	} else {
		switch packet.Command {
		case protocol.CMD_PING:
			response = "PONG"
			if len(packet.Key) > 0 {
				response += " " + packet.Key
			}
			break
		case protocol.CMD_SET:
			storage[packet.Key] = packet.Value
			fmt.Println(storage)
			response = "OK"
			break
		default:
			response = fmt.Sprintf("Invalid command %d", packet.Command)
		}
	}

	conn.Write([]byte(response + "\r\n"))
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
