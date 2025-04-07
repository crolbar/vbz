package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"vbz/audioCapture"
	"vbz/fft"
)

func (v *VBZ) parseEarlyArgs() error {
	for i, arg := range os.Args {
		switch arg {
		case "--decay":
			if i+1 >= len(os.Args) {
				return errors.New("params to --decay not enough, view --help")
			}

			num, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				return errors.New(fmt.Sprintf("error while parsing int at decay percentage argument"))
			}
			v.fft.Decay = num
		case "--filter-range", "-fr":
			if i+1 >= len(os.Args) {
				return errors.New("params to --filter-range/-fr not enough, view --help")
			}

			num, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				return errors.New(fmt.Sprintf("error while parsing int at filter range argument"))
			}
			v.fft.FilterRange = num
		case "--filter-mode", "-fm":
			if i+1 >= len(os.Args) {
				return errors.New("params to --filter-mode/-fm not enough, view --help")
			}
			switch strings.ToLower(os.Args[i+1]) {
			case "block":
				v.fft.FilterMode = fft.Block
			case "box", "boxfilter":
				v.fft.FilterMode = fft.BoxFilter
			case "dbox", "doublebox", "doubleboxfilter":
				v.fft.FilterMode = fft.DoubleBoxFilter
			case "none":
				v.fft.FilterMode = fft.None
			default:
				return errors.New("no such filter specifed to --filter-mode/-fm: " + os.Args[i+1])
			}

		case "--amp-scalar", "-as":
			if i+1 >= len(os.Args) {
				return errors.New("params to --amp-scalar/-as not enough, view --help")
			}

			num, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				return errors.New(fmt.Sprintf("error while parsing int at amp scalar argument"))
			}
			v.fft.AmpScalar = num
		case "--host":
			if i+1 >= len(os.Args) {
				return errors.New("params to --host not enough, view --help")
			}

			v.settings.Host = os.Args[i+1]
		case "--port":
			if i+1 >= len(os.Args) {
				return errors.New("params to --port not enough, view --help")
			}

			num, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				return errors.New(fmt.Sprintf("error while parsing int at port number argument"))
			}
			v.settings.Port = num
		case "-d", "--device-idx":
			if i+1 >= len(os.Args) {
				return errors.New("params to --device-idx/-d not enough, view --help")
			}

			num, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				return errors.New(fmt.Sprintf("error while parsing int at device index argument"))
			}
			v.settings.DeviceIdx = num
		case "-c", "--config":
			if i+1 >= len(os.Args) {
				return errors.New("params to --config/-c not enough, view --help")
			}

			v.configPath = os.Args[i+1]
		case "--debug":
			v.debug = true
		case "--fill-bins":
			v.fillBins = true
		case "--help", "-h":
			fmt.Println("TODO: help menu")
			v.shouldNotEnterTui = true
			return nil
		}

	}
	return nil
}

func (v *VBZ) parseLateArgs() error {
	for i, arg := range os.Args {
		switch arg {

		case "--off":
			v.turtOffRGB()
			v.shouldNotEnterTui = true
		case "--red":
			v.setAllLEDsToColor(255, 0, 0)
			v.shouldNotEnterTui = true
		case "--green":
			v.setAllLEDsToColor(0, 255, 0)
			v.shouldNotEnterTui = true
		case "--blue":
			v.setAllLEDsToColor(0, 0, 255)
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

			err = v.setAllLEDsToColor(uint8(r), uint8(g), uint8(b))
			if err != nil {
				return err
			}

			v.shouldNotEnterTui = true
		}
	}
	return nil
}
