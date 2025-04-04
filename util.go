package main

import (
	"bytes"
	"encoding/binary"
)

func byteToU8(data []byte) ([]uint8, error) {
	buf := bytes.NewReader(data)

	var result []uint8

	for {
		var u uint8
		err := binary.Read(buf, binary.LittleEndian, &u)
		if err != nil {
			break
		}
		result = append(result, uint8(u))
	}

	return result, nil
}
