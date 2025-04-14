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

func (s *Settings) ParseEarlyArgs() error {
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
			panic("TODO HELP MENU")
			// v.shouldNotEnterTui = true
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
			return err
		}
	}

	return nil
}
