package main

import (
	"bytes"
	"encoding/binary"
)

func byteToFloat32(data []byte) ([]float32, error) {
	buf := bytes.NewReader(data)

	var result []float32

	for {
		var f float32
		err := binary.Read(buf, binary.LittleEndian, &f)
		if err != nil {
			break
		}
		result = append(result, f)
	}

	return result, nil
}
