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
			readMetadata(reader)
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
			//packet 3
			packetLen, packetID, err := readMetadata(reader)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Len", packetLen, "packetID", packetID)
			packetLen, packetID, err = readMetadata(reader)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(packetLen, packetID)
			buf := make([]byte, 100)
			reader.Read(buf)
			fmt.Println("newBuf", buf)
			sendConf(conn)

			readConfig(reader)
	
			sendConfConfirmation(conn)
			buf = make([]byte, 100)
			reader.Read(buf)
			fmt.Println("newBuf", buf)
		} else if nextState == 1 {
			readMetadata(reader)

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
	_, err := reader.Read(buf)
	if err != nil {
		fmt.Println("DSADSADSA", err)
	}
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

func sendLoginSuccess(conn net.Conn, uuid string, username string) {
	writer := bufio.NewWriter(conn)

	converted, err := ConvertUUID(uuid)
	if err != nil {
		fmt.Println(err)
	}
	packetLen := len(converted.Bytes()) + len([]byte(username)) + 4

	var buffer bytes.Buffer
	WriteVarInt(&buffer, packetLen)
	WriteVarInt(&buffer, 0x02)

	buffer.Write([]byte(converted.Bytes()))
	WriteVarInt(&buffer, len(username))
	buffer.Write([]byte(username))
	WriteVarInt(&buffer, 0x00)
	WriteVarInt(&buffer, 0x00)

	_, err = writer.Write(buffer.Bytes())
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

func sendConf(conn net.Conn) {

	writer := bufio.NewWriter(conn)

	var buffer bytes.Buffer
	id := 0x0E
	packCount := 1
	nameSpace := "minecraft"
	nameSpaceID := "core"
	version := "1.21"

	packetLen := len(nameSpace) + len(nameSpaceID) + len(version) + 5

	WriteVarInt(&buffer, packetLen)
	WriteVarInt(&buffer, id)

	WriteVarInt(&buffer, packCount)
	WriteVarInt(&buffer, len(nameSpace))
	buffer.Write([]byte(nameSpace))

	WriteVarInt(&buffer, len(nameSpaceID))
	buffer.Write([]byte(nameSpaceID))

	WriteVarInt(&buffer, len(version))
	buffer.Write([]byte(version))

	fmt.Println("Packetlen", packetLen)
	fmt.Println("Packetlen SENT", len(buffer.Bytes()))

	fmt.Println("BUFFFFFFFER", buffer.Bytes())

	_, err := writer.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error sending config packet:", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}

	fmt.Println("Packet sent confing")
}

func sendConfConfirmation(conn net.Conn) {
	writer := bufio.NewWriter(conn)

	var buffer bytes.Buffer
	WriteVarInt(&buffer, 0x01)
	WriteVarInt(&buffer, 0x03)

	_, err := writer.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error sending config packet:", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}

	fmt.Println("Packet confirmation sent")
}




// func login(conn net.Conn) {

// 	var buffer bytes.Buffer
// 	writer := bufio.NewWriter(conn)

// 	entityID := 1
// 	isHardCore := 0

// 	buffer.WriteByte(byte(entityID))
// 	err := binary.Write(&buffer, binary.BigEndian, entityID)

// }
