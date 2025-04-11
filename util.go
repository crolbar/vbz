package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"strconv"
	"time"
	"vbz/audioCapture"
)

type Refresh struct{}

func (v *VBZ) triggerRefresh() {
	if v.p == nil {
		return
	}
	v.p.Send(Refresh{})
}
func (v *VBZ) triggerLaterRefresh() {
	go func() {
		time.Sleep(time.Millisecond * 16)
		if v.p == nil {
			return
		}
		v.p.Send(Refresh{})
	}()
}

func (v *VBZ) parseLateArgs() error {
	for i, arg := range os.Args {
		switch arg {

		case "--off":
			v.led.TurtOffRGB()
			v.shouldNotEnterTui = true
		case "--red":
			v.led.SetAllLEDsToColor(255, 0, 0)
			v.shouldNotEnterTui = true
		case "--green":
			v.led.SetAllLEDsToColor(0, 255, 0)
			v.shouldNotEnterTui = true
		case "--blue":
			v.led.SetAllLEDsToColor(0, 0, 255)
			v.shouldNotEnterTui = true

		case "--list-devices", "-l":
			audioCapture.PrintDevices()
			v.shouldNotEnterTui = true

		case "--set-color":
			if i+3 >= len(os.Args) {
				return errors.New("params to --set-color not enough, view --help")
			}

			r, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				return err
			}
			g, err := strconv.Atoi(os.Args[i+2])
			if err != nil {
				return err
			}
			b, err := strconv.Atoi(os.Args[i+3])
			if err != nil {
				return err
			}

			err = v.led.SetAllLEDsToColor(uint8(r), uint8(g), uint8(b))
			if err != nil {
				return err
			}

			v.shouldNotEnterTui = true
		}
	}
	return nil
}

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
