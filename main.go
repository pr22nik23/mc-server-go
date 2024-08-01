package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type Server struct {
}

const (
	// Default Minecraft server port
	port = ":25565"
)

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection from", conn.RemoteAddr())
	reader := bufio.NewReader(conn)

	// Handle handshake
	packetLength, err := binary.ReadUvarint(reader)
	if err != nil {
		switch {
		case errors.Is(err, io.EOF):
			fmt.Println("No bytes were read", err)
		case errors.Is(err, io.ErrUnexpectedEOF):
			fmt.Println("Unexpteced err", err)
		default:
			fmt.Println("Error:", err)
		}
	}
	fmt.Println("Packet length:", packetLength)
	packetID, err := binary.ReadUvarint(reader)
	if err != nil {
		switch {
		case errors.Is(err, io.EOF):
			fmt.Println("No bytes were read", err)
		case errors.Is(err, io.ErrUnexpectedEOF):
			fmt.Println("Unexpteced err", err)
		default:
			fmt.Println("Error:", err)
		}
	}

	if packetID == 0 {
		fmt.Println("Starting packet")
		protocol_version, _ := binary.ReadUvarint(reader)
		fmt.Println("protocol", protocol_version)
		server_aadress, _ := ReadVarString(reader)
		fmt.Println("server address", server_aadress)
		buf := make([]byte, 2)
		_, err := reader.Read(buf)
		if err != nil {
			fmt.Println("Err", err)
		}

		port := binary.ByteOrder.Uint16(binary.BigEndian, buf)
		fmt.Println("Port", port)

		nextState, err := binary.ReadUvarint(reader)
		fmt.Println("Nexct state", nextState)
		ed, err := binary.ReadUvarint(reader)
		if err != nil {
			fmt.Println("Yep", err)
		}
		fmt.Println(ed)
		// fmt.Println("Read packages", read)
	} else if packetID == 122 {
		fmt.Println("somethign else ")
	}
}

func ReadVarString(reader *bufio.Reader) (string, error) {
	// Read the length of the string as a varint.
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read string length: %v", err)
	}

	fmt.Println("string len", length)
	// Allocate a buffer to hold the string bytes.
	buf := make([]byte, length)

	// Read the bytes of the string into the buffer.
	_, err = reader.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read string data: %v", err)
	}

	// Convert the byte slice to a string and return it.
	return string(buf), nil
}
