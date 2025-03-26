package orgb

import (
	"bytes"
	bn "encoding/binary"
	"fmt"
	"net"
)

const HEADER_SIZE int = 4 * 4

func Connect(host string, port int) (*ORGBConn, error) {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	conn := ORGBConn{conn: c}

	name := []byte("vbz")
	err = conn.sendMessage(NET_PACKET_ID_SET_CLIENT_NAME, 0, &name)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func (c *ORGBConn) Close() error {
	return c.conn.Close()
}

func (c *ORGBConn) GetControllerCount() (int, error) {
	err := c.sendMessage(NET_PACKET_ID_REQUEST_CONTROLLER_COUNT, 0, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.readMessage()
	if err != nil {
		return 0, err
	}

	n := int(bn.LittleEndian.Uint32(resp))

	return n, nil
}

func (c *ORGBConn) GetController(devId int) (Controller, error) {
	err := c.sendMessage(NET_PACKET_ID_REQUEST_CONTROLLER_DATA, devId, nil)
	if err != nil {
		return Controller{}, err
	}

	buff, err := c.readMessage()
	if err != nil {
		return Controller{}, err
	}

	controller, err := parseControllerData(buff)

	return controller, nil
}

func (c *ORGBConn) UpdateLEDS(devId int, colors []RGBColor) error {
	var (
		numColors = len(colors)
		size      = 2 + // num colors
			(4 * numColors) // colors

		colorsBuf bytes.Buffer
	)
	colorsBuf.Grow(size + 4)

	// data size
	{
		sizeBuf := make([]byte, 4)
		bn.LittleEndian.PutUint32(sizeBuf, uint32(size))
		colorsBuf.Write(sizeBuf)
	}

	// num colors
	{
		numBuf := make([]byte, 2)
		bn.LittleEndian.PutUint16(numBuf, uint16(numColors))
		colorsBuf.Write(numBuf)
	}

	// colors
	{
		for _, c := range colors {
			colorsBuf.WriteByte(c.Red)
			colorsBuf.WriteByte(c.Green)
			colorsBuf.WriteByte(c.Blue)
			colorsBuf.WriteByte(0)
		}
	}

	buf := colorsBuf.Bytes()
	return c.sendMessage(NET_PACKET_ID_RGBCONTROLLER_UPDATELEDS, devId, &buf)
}

func (c *ORGBConn) sendMessage(command, deviceID int, buffer *[]byte) error {
	size := 0
	if buffer != nil {
		size = len(*buffer)
	}

	header := encodeHeader(ORGBHeader{
		devIdx: uint32(deviceID),
		id:     uint32(command),
		size:   uint32(size),
	})

	if buffer != nil {
		header.Write(*buffer)
	}

	_, err := c.conn.Write(header.Bytes())

	return err
}

func (c *ORGBConn) readMessage() ([]byte, error) {
	buf := make([]byte, HEADER_SIZE)
	_, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	header := decodeHeader(buf)
	buf = make([]byte, header.size)
	_, err = c.conn.Read(buf)

	return buf, err
}
