package orgb

import (
	bn "encoding/binary"
	"fmt"
)

func parseControllerData(data []byte) (Controller, error) {
	var (
		c   = Controller{}
		off = 4 // data size
	)

	c.Type = int(bn.LittleEndian.Uint32(data[off:]))
	off += 4 // type

	// Name
	{
		size := bn.LittleEndian.Uint16(data[off:])
		off += 2 // name_size

		c.Name = string(data[off : uint16(off)+size-1])
		off += int(size) // name
	}

	// Description
	{
		size := bn.LittleEndian.Uint16(data[off:])
		off += 2 // desc_size

		c.Description = string(data[off : uint16(off)+size-1])
		off += int(size) // desc
	}

	// Version
	{
		size := bn.LittleEndian.Uint16(data[off:])
		off += 2 // version_size

		c.Version = string(data[off : uint16(off)+size-1])
		off += int(size) // version
	}

	// Serial
	{
		size := bn.LittleEndian.Uint16(data[off:])
		off += 2 // serial_size

		c.Serial = string(data[off : uint16(off)+size-1])
		off += int(size) // serial
	}

	// Location
	{
		size := bn.LittleEndian.Uint16(data[off:])
		off += 2

		c.Location = string(data[off : uint16(off)+size-1])
		off += int(size)
	}

	// modes
	{
		numModes := bn.LittleEndian.Uint16(data[off:])
		off += 2 // num_modes

		c.ActiveMode = int(bn.LittleEndian.Uint32(data[off:]))
		off += 4 // active mode

		off += getModesSize(int(numModes), data[off:]) // skip modes
	}

	// zones
	{
		numZones := bn.LittleEndian.Uint16(data[off:])
		off += 2 // num zones

		off += getZonesSize(int(numZones), data[off:]) // skip zones
	}

	// LEDS
	{
		numLEDS := bn.LittleEndian.Uint16(data[off:])
		off += 2 // num zones

		off += getLEDsSize(int(numLEDS), data[off:]) // skip leds
	}

	// Colors
	{
		numColors := bn.LittleEndian.Uint16(data[off:])
		off += 2 // num colors

		colors := make([]RGBColor, numColors)

		for i := range numColors {
			colors[i] = RGBColor{
				Red:   data[off],
				Green: data[off+1],
				Blue:  data[off+2],
			}
			off += 4 * 1
		}

		c.Colors = colors
	}

	return c, nil
}

func (c Controller) String() string {
	return fmt.Sprintf(
		"Type: %d,\nName: %s\nDesc: %s\nVersion: %s, Serial: %s, Location: %s, ActiveMode: %d\nColors: %v",
		c.Type,
		c.Name,
		c.Description,
		c.Version,
		c.Serial,
		c.Location,
		c.ActiveMode,
		c.Colors,
	)
}
