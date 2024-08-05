package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type Server struct {
}

const (
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

	reader := bufio.NewReader(conn)
	fmt.Println("New connection from", conn.RemoteAddr())
	// buf := make([]byte, 1024)
	// reader.Read(buf)
	// fmt.Println("This is hte one", buf)
	packetLen, packetID, err := readMetadata(reader)
	switch {
	case errors.Is(err, io.EOF):
		buf := make([]byte, 100)
		reader.Read(buf)
		fmt.Println("buff", buf)
		fmt.Println("EOF", err)
		return
	case errors.Is(err, io.ErrUnexpectedEOF):
		fmt.Println("ErrUnexpectedEOF", err)
		return
	}
	fmt.Printf("Packet Length: %d, PacketID: %d\n", packetLen, packetID)
	if packetID == 0 {
		protocol_version, _ := binary.ReadUvarint(reader)
		fmt.Println("Protocol: ", protocol_version)
		server_aadress, _ := ReadVarString(reader)
		fmt.Println("Server Address: ", server_aadress)
		buf := make([]byte, 2)
		_, err := reader.Read(buf)
		if err != nil {
			fmt.Println("Err", err)
		}

		port := binary.ByteOrder.Uint16(binary.BigEndian, buf)
		fmt.Println("Port: ", port)

		nextState, err := binary.ReadUvarint(reader)
		if err != nil {
			fmt.Println("Erroorr", err)
		}

		if nextState == 2 {
			handleLogin(reader)
			username, err := ReadVarString(reader)
			if err != nil {
				fmt.Println("Error", err)
			}
			fmt.Println("Username", username)

			uuid, err := ReadUUID(reader)
			if err != nil {
				fmt.Println("Error", err)
			}
			fmt.Println("UUID", uuid)

			sendLoginSuccess(conn, uuid, username)
			handleLogin(reader)

		} else if nextState == 1 {
			handleLogin(reader)

			SendStatusResponse(conn)

		} else {
			fmt.Println("Some unknown state")
			return
		}

	} else if packetID == 122 {
		fmt.Println("We are in 122")
		readReportDetails(reader)
	} else {
		fmt.Println("This isjust a test fir sinetubg ekse:")
	}
}

func readReportDetails(reader *bufio.Reader) {
	buf := make([]byte, 255)
	read, err := reader.Read(buf)
	if err != nil {
		fmt.Println("DSADSADSA", err)
	}
	fmt.Println("Read bytes", read)
	fmt.Println(buf)
	fmt.Println(string(buf))
	count, err := binary.ReadVarint(reader)
	if err != nil {
		fmt.Println("Read Error", err)
	}

	fmt.Println("Count", count)
	for i := 0; int64(i) < count; i++ {
		fmt.Println("I", i)
		title, err := ReadVarString(reader)
		if err != nil {
			fmt.Println("String error", err)
		}
		description, err := ReadVarString(reader)
		if err != nil {
			fmt.Println("String error", err)
		}

		fmt.Println("Title -", title, "Description -", description)
	}

}

func readMetadata(reader *bufio.Reader) (uint64, uint64, error) {
	packetLength, err := binary.ReadUvarint(reader)
	if err != nil {
		return 0, 0, err
	}
	packetID, err := binary.ReadUvarint(reader)
	if err != nil {
		return 0, 0, err
	}

	return packetLength, packetID, nil
}

func handleLogin(reader *bufio.Reader) {
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("length", length)

	packet_id, err := binary.ReadUvarint(reader)
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("new packet id", packet_id)
}

func sendLoginSuccess(conn net.Conn, uuid string, username string) {
	writer := bufio.NewWriter(conn)

	packetLen := len([]byte(uuid)) + len([]byte(username)) + 2

	var buffer bytes.Buffer
	WriteVarInt(&buffer, packetLen)
	WriteVarInt(&buffer, 0)

	fmt.Println("Packetlen", packetLen)
	fmt.Println(buffer.Bytes())

	// packet := []byte()
	// buffer.Write([]byte())
	buffer.Write([]byte(uuid))
	buffer.Write([]byte(username))
	// packet = append(packet, packetID)
	// packet = append(packet, []byte(uuid)...)
	// packet = append(packet, []byte(username)...)
	// packet = append(packet, buf...)

	fmt.Println("package we are sending", buffer.Bytes())
	_, err := writer.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error sending login success packet:", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
	fmt.Println("Login success packet sent")
}

func handleReadError(err error) {
	switch {
	case errors.Is(err, io.EOF):
		fmt.Println("No bytes were read:", err)
	case errors.Is(err, io.ErrUnexpectedEOF):
		fmt.Println("Unexpected EOF:", err)
	default:
		fmt.Println("Error:", err)
	}
}
