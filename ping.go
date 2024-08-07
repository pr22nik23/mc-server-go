package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"

	"github.com/gofrs/uuid/v5"
)

func SendStatusResponse(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	var buffer bytes.Buffer

	uuid, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Server error:", err)
		return
	}

	statusResponse := StatusResponse{
		Version: Version{
			Protocol: 767,
			Name:     "1.21",
		},
		Players: Players{
			Max:    100,
			Online: 5,
			Sample: []Sample{
				{
					Name: "tere",
					Id:   uuid,
				},
			},
		},
		Description: Description{
			Text: "Harri on winners attitude",
		},
		EnforcesSecureChat: false,
	}

	res, err := json.Marshal(statusResponse)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}

	packetLength := len(res) + 3

	WriteVarInt(&buffer, packetLength)
	WriteVarInt(&buffer, 0x00)
	WriteVarInt(&buffer, len(res))
	fmt.Println("packet Len ", packetLength)
	fmt.Println("Actual packet len", len(buffer.Bytes()))
	buffer.Write(res)

	_, err = writer.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error sending status response packet:", err)
		return
	}

	// Flush the writer to ensure data is sent.
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}

	fmt.Println("Sent status response")
}
