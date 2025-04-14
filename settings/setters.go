package settings

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	ft "vbz/fft/filter_types"
)

type setIntConfigType = func(field *int, val string, lineNumber int) error

func setIntConfig(field *int, val string, lineNumber int) error {
	num, err := strconv.Atoi(val)
	if err != nil {
		return errors.New(fmt.Sprintf("error while parsing int at line: %d", lineNumber))
	}

	*field = num
	return nil
}

type setIntArgsType = func(field *int, argIdx int, argKey string) error

func setIntArgs(field *int, argIdx int, argKey string) error {
	if argIdx+1 >= len(os.Args) {
		return errors.New(
			fmt.Sprintf("params to %s not enough, view --help", argKey),
		)
	}

	num, err := strconv.Atoi(os.Args[argIdx+1])
	if err != nil {
		return errors.New(fmt.Sprintf("error while parsing int at argument: %s", argKey))
	}

	*field = num
	return nil
}

type setStringConfigType = func(field *string, val string)

func setStringConfig(field *string, val string) {
	*field = val
}

type setStringArgsType = func(field *string, argIdx int, argKey string) error

func setStringArgs(field *string, argIdx int, argKey string) error {
	if argIdx+1 >= len(os.Args) {
		return errors.New(
			fmt.Sprintf("params to %s not enough, view --help", argKey),
		)
	}

	*field = os.Args[argIdx+1]
	return nil
}

type setBoolConfigType = func(field *bool, val string, lineNumber int) error

func setBoolConfig(field *bool, val string, lineNumber int) error {
	switch val {
	case "true":
		*field = true
	case "false":
		*field = false
	default:
		return errors.New(fmt.Sprintf(
			"incorrect boolean value at line: %d", lineNumber),
		)
	}
	return nil
}

type setBoolTrueArgsType = func(field *bool)

func setBoolTrueArgs(field *bool) {
	*field = true
}

type setBoolFalseArgsType func(field *bool)

func setBoolFalseArgs(field *bool) {
	*field = false
}

type setFloatConfigType = func(field *float64, val string, lineNumber int) error

func setFloatConfig(field *float64, val string, lineNumber int) error {
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("error while parsing float64 at line: %d", lineNumber))
	}

	*field = num
	return nil
}

type setFloatArgsType = func(field *float64, argIdx int, argKey string) error

func setFloatArgs(field *float64, argIdx int, argKey string) error {
	if argIdx+1 >= len(os.Args) {
		return errors.New(
			fmt.Sprintf("params to %s not enough, view --help", argKey),
		)
	}

	num, err := strconv.ParseFloat(os.Args[argIdx+1], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("error while parsing float64 at argument: %s", argKey))
	}

	*field = num
	return nil
}

type setFilterModeConfigType = func(field *ft.FilterType, val string, lineNumber int) error

func setFilterModeConfig(field *ft.FilterType, val string, lineNumber int) error {
	switch val {
	case "None":
		*field = ft.None
	case "Block":
		*field = ft.Block
	case "BoxFilter":
		*field = ft.BoxFilter
	case "DoubleBoxFilter":
		*field = ft.DoubleBoxFilter
	default:
		return errors.New(fmt.Sprintf(
			"incorrect filter type at line: %d (see /fft/filter_types.go for types)", lineNumber),
		)
	}
	return nil
}

type setFilterModeArgsType = func(field *ft.FilterType, argIdx int, argKey string) error

func setFilterModeArgs(field *ft.FilterType, argIdx int, argKey string) error {
	if argIdx+1 >= len(os.Args) {
		return errors.New(
			fmt.Sprintf("params to %s not enough, view --help", argKey),
		)
	}

	switch os.Args[argIdx+1] {
	case "None":
		*field = ft.None
	case "Block":
		*field = ft.Block
	case "BoxFilter":
		*field = ft.BoxFilter
	case "DoubleBoxFilter":
		*field = ft.DoubleBoxFilter
	default:
		return errors.New(fmt.Sprintf(
			"incorrect filter type at for arg: %s (see /fft/filter_types.go for types)", argKey),
		)
	}
	return nil
}
