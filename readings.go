package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
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

func ConvertUUID(uuid string) (*big.Int, error) {
	cleanUUID := strings.ReplaceAll(uuid, "-", "")

	intValue, success := new(big.Int).SetString(cleanUUID, 16)
	if !success {
		return nil, fmt.Errorf("error converting UUID to integer")
	}

	return intValue, nil
}

func readNBT(r *bufio.Reader, input interface{}) error {

	decoder := json.NewDecoder(r)

	err := decoder.Decode(input)
	if err != nil {
		return err
	}

	return nil
}

func readConfig(reader *bufio.Reader) {

	_, packetID, err := readMetadata(reader)
	if err == io.EOF {
		if errors.Is(err, io.EOF) {
			fmt.Println(err)
			return
		}
	}

	if packetID == 7 {
		count, err := binary.ReadUvarint(reader)
		if err != nil {
			fmt.Println(err)
		}

		namespace, err := ReadVarString(reader)
		if err != nil {
			fmt.Println(err)
		}

		ID, err := ReadVarString(reader)
		if err != nil {
			fmt.Println(err)
		}

		version, err := ReadVarString(reader)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Count", count)
		fmt.Println("Namespace", namespace)
		fmt.Println("ID", ID)
		fmt.Println("Version", version)

	} else {
		packetLen, packetID, err := readMetadata(reader)
		if err != nil {
			fmt.Println(err)
		}

		buf := make([]byte, packetLen)
		reader.Read(buf)

		fmt.Printf("Optional Buffer{ ID: %d, Buffer: %v}", packetID, buf)

		readConfig(reader)
	}
}
