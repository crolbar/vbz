package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"vbz/audioCapture"
)

func getConfigPathFromArgs() (string, error) {
	for i, arg := range os.Args {
		switch arg {
		case "-c", "--config":
			if i+1 >= len(os.Args) {
				return "", errors.New("params to --config/-c not enough, view --help")
			}

			return os.Args[i+1], nil
		}
	}
	return "", nil
}

func (v *VBZ) parseArgs() (bool, error) {
	for i, arg := range os.Args {
		switch arg {

		case "--off":
			v.turtOffRGB()
			return true, nil
		case "--red":
			v.setAllLEDsToColor(255, 0, 0)
			return true, nil
		case "--green":
			v.setAllLEDsToColor(0, 255, 0)
			return true, nil
		case "--blue":
			v.setAllLEDsToColor(0, 0, 255)
			return true, nil

        case "--list-devices", "-l":
			audioCapture.PrintDevices()
			return true, nil

		case "--set-color":
			if i+3 >= len(os.Args) {
				return false, errors.New("params to --set-color not enough, view --help")
			}

			r, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				return false, err
			}
			g, err := strconv.Atoi(os.Args[i+2])
			if err != nil {
				return false, err
			}
			b, err := strconv.Atoi(os.Args[i+3])
			if err != nil {
				return false, err
			}

			err = v.setAllLEDsToColor(uint8(r), uint8(g), uint8(b))
			if err != nil {
				return false, err
			}

			return true, nil

		case "--help", "-h":
			fmt.Println("TODO: help menu")
			return true, nil
		}
	}

	return false, nil
}
