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

	reader := bufio.NewReader(conn)
	fmt.Println("New connection from", conn.RemoteAddr())
	for {
		_, packetID, err := readMetadata(reader)
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
		fmt.Println("packet", packetID)
		if packetID == 0 {
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
			if err != nil {
				fmt.Println("TEST", err)
			}

			fmt.Println("Nexct state", nextState)

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
				// packet_id, err := binary.ReadUvarint(reader)
				// if err != nil {
				// 	fmt.Println("err", err)
				// }

				// fmt.Println("Newpacked", packet_id)
			}

		} else if packetID == 122 {
			fmt.Println("somethign else ")
		} else {
			fmt.Println("This isjust a test fir sinetubg ekse:")
		}
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
	fmt.Println([]byte(uuid))
	writer := bufio.NewWriter(conn)
	buf := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutVarint(buf, 0x00)
	packetID := byte(0x02)
	packetLen := len([]byte(uuid)) + len([]byte(username)) + 2
	packet := []byte{byte(packetLen)}
	packet = append(packet, packetID)
	packet = append(packet, []byte(uuid)...)
	packet = append(packet, []byte(username)...)
	// packet = append(packet, buf...)

	_, err := writer.Write(packet)
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

func ReadVarString(reader *bufio.Reader) (string, error) {
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read string length: %v", err)
	}

	buf := make([]byte, length)
	_, err = reader.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read string data: %v", err)
	}

	return string(buf), nil
}

func ReadUUID(reader *bufio.Reader) (string, error) {
	uuidBytes := make([]byte, 16)
	_, err := io.ReadFull(reader, uuidBytes)
	if err != nil {
		return "", fmt.Errorf("failed to read UUID: %v", err)
	}

	// Format UUID with dashes
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		uuidBytes[0:4],
		uuidBytes[4:6],
		uuidBytes[6:8],
		uuidBytes[8:10],
		uuidBytes[10:16])

	return uuid, nil
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
