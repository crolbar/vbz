package settings

import (
	"errors"
	"fmt"
	"os"
	ft "vbz/fft/filter_types"
)

func parseConfigPathArg() (path string, err error) {
	for i := 0; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--config", "-c":
			if i+1 >= len(os.Args) {
				err = errors.New(
					fmt.Sprintf("params to %s not enough, view --help", os.Args[i]),
				)
				break
			}

			path = os.Args[i+1]
			break
		}
	}

	return
}

func (s *Settings) ParseEarlyArgs() (bool, error) {
	var (
		funcType interface{}
		ok       bool
		skipNext = false

		key string
		err error
	)

	for i, arg := range os.Args[1:] {
		if skipNext {
			skipNext = false
			continue
		}

		switch arg {
		case "--help", "-h":
			fmt.Println(
				"VBZ - OpenRGB Audio Visualizer Client" + "\n" +
					"Usage: vbz [OPTION..]" + "\n" +
					"\n" + "[OPTIONS]" + "\n" +
					"-d, --device-idx    Device to capture (use -l to see devices)" + "\n" +
					"-c, --config        Config file path" + "\n" +
					"--host              OpenRGB server host" + "\n" +
					"--port              OpenRGB server port" + "\n" +
					"--no-leds           Don't update leds" + "\n" +
					"--no-open-rgb       Don't try to connect to OpenRGB server" + "\n" +
					"--fill-bins         Fills side black bars in the bins visualizer" + "\n" +
					"--hue-rate          Hue change rate (speeds up color wave)" + "\n" +
					"--amp-scalar        scale up the amplitude for better visualization" + "\n" +
					"--filter-mode       Averaging filter type" + "\n" +
					"--filter-range      Range of the averaging filter" + "\n" +
					"--decay             Percentage of decay of amplitude in each frame",
			)
			return true, nil
		}

		if funcType, ok, key = getFuncType(arg); !ok {
			continue
		}

		key = kebabToPascal(key)

		switch f := funcType.(type) {
		case setIntArgsType:
			err = f(getFieldPointer(s, key).(*int), i+1, arg) // the +1 is because we skip the first arg in the loop
			skipNext = true
		case setStringArgsType:
			f(getFieldPointer(s, key).(*string), i+1, arg)
			skipNext = true
		case setBoolTrueArgsType:
			f(getFieldPointer(s, key).(*bool))
		case setBoolFalseArgsType:
			f(getFieldPointer(s, key).(*bool))
		case setFloatArgsType:
			err = f(getFieldPointer(s, key).(*float64), i+1, arg)
			skipNext = true
		case setFilterModeArgsType:
			err = f(getFieldPointer(s, key).(*ft.FilterType), i+1, arg)
			skipNext = true
		}

		if err != nil {
			return false, err
		}
	}

	return false, nil
}
