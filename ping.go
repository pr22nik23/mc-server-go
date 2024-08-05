package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"github.com/gofrs/uuid/v5"
)

func SendStatusResponse(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	packetID := byte(0x00)

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

	varString, err := EncodeStringWithVarInt(string(res))
	if err != nil {
		fmt.Println(err)
	}
	packetLength := 3 + len(varString)
	packetLengthVarInt := VarInt(int64(packetLength))

	// Construct the final packet.
	packet := append(packetLengthVarInt, packetID)
	packet = append(packet, varString...)

	fmt.Println("Packet len", len(packet))
	fmt.Println("Written packet len", packetLength)
	// Write the packet to the connection.
	_, err = writer.Write(packet)
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
