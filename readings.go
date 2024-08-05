package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

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

func EncodeStringWithVarInt(s string) ([]byte, error) {
	// Get the string representation of the struct

	// Calculate the number of UTF-16 code units
	// utf16Length := CalculateUTF16Length(str)
	// if utf16Length > MaxUTF16Units {
	// 	return nil, fmt.Errorf("string exceeds maximum length in UTF-16 code units: %d", utf16Length)
	// }

	// Calculate the number of UTF-8 bytes
	utf8Length := len([]byte(s))
	// if utf8Length > MaxUTF16Units*3 {
	// 	return nil, fmt.Errorf("string exceeds maximum byte length: %d", utf8Length)
	// }

	// Create a buffer to write the VarInt and UTF-8 string
	var buffer bytes.Buffer

	// Write the length as a VarInt
	lenBytes := VarInt(int64(utf8Length))
	buffer.Write(lenBytes)

	// Write the UTF-8 string bytes
	buffer.WriteString(s)

	return buffer.Bytes(), nil
}


func VarInt(number int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	size := binary.PutVarint(buf, number)

	return buf[:size]
}


func WriteVarInt(buffer *bytes.Buffer, value int) {
	for {
		temp := byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			temp |= 0x80
		}
		buffer.WriteByte(temp)
		if value == 0 {
			break
		}
	}
}