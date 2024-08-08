package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

type State int

const (
	Handshake State = iota
	Status
	Login
	Transfer
)

type Client struct {
	Conn  net.Conn
	State State
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn:  conn,
		State: Handshake,
	}
}

func (c *Client) WriteToConn(msg []byte) {
	writer := bufio.NewWriter(c.Conn)
	_, err := writer.Write(msg)
	if err != nil {
		fmt.Println("Error sending packet: ", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
	fmt.Println("Package successfully sent")
}

func (c *Client) ReadLoop() {
	defer c.Conn.Close()
	reader := bufio.NewReader(c.Conn)
	fmt.Println("New client connected", c.Conn.LocalAddr())
	for {
		packetLen, packetID, err := readMetadata(reader)
		if err != nil {
			fmt.Println("Error reading packet's metadata")
		}
		switch c.State {
		case Handshake:
			fmt.Println("Handshaking request")
			c.ResolveHandshake(reader)
		case Status:
			fmt.Println("Ping request")
		case Login:
			fmt.Println("Login state")
		case Transfer:
			fmt.Println("Transfer state")
		}

		fmt.Println(packetLen, packetID)
	}
}

func (c *Client) ResolveHandshake(reader *bufio.Reader) {
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

	c.State = State(int(nextState))
}
