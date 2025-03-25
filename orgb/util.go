package orgb

import (
	"bytes"
	"encoding/binary"
	bn "encoding/binary"
)

func encodeHeader(header ORGBHeader) *bytes.Buffer {
	b := bytes.NewBufferString("ORGB")

	for _, v := range []uint32{
		header.devIdx,
		header.id,
		header.size,
	} {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, v)
		b.Write(buf)
	}

	return b
}

func decodeHeader(buffer []byte) ORGBHeader {
	return ORGBHeader{
		binary.LittleEndian.Uint32(buffer[4:]),
		binary.LittleEndian.Uint32(buffer[8:]),
		binary.LittleEndian.Uint32(buffer[12:]),
	}
}

func getLEDsSize(numLEDs int, data []byte) int {
	off := 0

	for range numLEDs {
		nameLen := bn.LittleEndian.Uint16(data[off:])
		off += 2            // name len
		off += int(nameLen) // name
		off += 4            // value
	}

	return off
}

func getZonesSize(numZones int, data []byte) int {
	off := 0

	for range numZones {
		nameLen := bn.LittleEndian.Uint16(data[off:])
		off += 2            // name len
		off += int(nameLen) // name

		off += 4 + // type
			4 + // leds min
			4 + // leds max
			4 // leds count

		matrixLen := bn.LittleEndian.Uint16(data[off:])
		off += 2              // matrix len
		off += int(matrixLen) // matrix
	}

	return off
}

func getModesSize(numModes int, data []byte) int {
	off := 0

	for range numModes {
		nameLen := bn.LittleEndian.Uint16(data[off:])
		off += 2 // name len

		// fmt.Println(nameLen)
		// fmt.Println("mode name: ", string(data[off:off+int(nameLen)-1]))

		off += int(nameLen) // name

		off += 4 + // value
			4 + // flags
			4 + // speed min
			4 + // speed max
			4 + // colors min
			4 + // colors max
			4 + // speed
			4 + // direction
			4 // color mode

		colorsNum := bn.LittleEndian.Uint16(data[off:])
		off += 2 // colors num

		off += int(colorsNum) * (4 * 1) // colors
	}

	return off
}
