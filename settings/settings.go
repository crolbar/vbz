package settings

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	ft "vbz/fft/filter_types"
)

// in config file keys are case insensitive
type Settings struct {
	Port        int
	Host        string
	DeviceIdx   int
	Debug       bool
	FillBins    bool
	HueRate     float64
	AmpScalar   int
	FilterMode  ft.FilterType
	FilterRange int
	Decay       int
	Config      string // config path
}

var DefaultSettings Settings = Settings{
	Port:      6742,
	Host:      "localhost",
	DeviceIdx: 0,
	Debug:     false,
	FillBins:  false,
	HueRate:   0.003 * 3, // 0.003 is 1 degree a tick

	AmpScalar:   5000,
	FilterRange: 2,
	FilterMode:  ft.DoubleBoxFilter,
	Decay:       80,
}

var fieldMapping = map[string]interface{}{
	// config file keys (see ./util.go for aliases for these)
	"Port":        setIntConfig,
	"Host":        setStringConfig,
	"DeviceIdx":   setIntConfig,
	"FillBins":    setBoolConfig,
	"HueRate":     setFloatConfig,
	"AmpScalar":   setIntConfig,
	"FilterRange": setIntConfig,
	"FilterMode":  setFilterModeConfig,
	"Decay":       setIntConfig,

	// cli args
	"--config":       setStringArgs,
	"--device-idx":   setIntArgs,
	"--host":         setStringArgs,
	"--port":         setIntArgs,
	"--debug":        setBoolTrueArgs,
	"--fill-bins":    setBoolTrueArgs,
	"--hue-rate":     setFloatArgs,
	"--amp-scalar":   setIntArgs,
	"--filter-range": setIntArgs,
	"--filter-mode":  setFilterModeArgs,
	"--decay":        setIntArgs,
}

func (s *Settings) InitSettings() error {
	var err error
	if s.Config == "" {
		err = s.getSettingsDefaultPath()
	} else {
		err = s.getSettings(s.Config)
	}
	if err != nil {
		return err
	}

	// VERY BAD SOLUTION FIX THIS !!
	// set uninited values to default
	if s.DeviceIdx == -1 {
		s.DeviceIdx = DefaultSettings.DeviceIdx
	}
	if s.Port == -1 {
		s.Port = DefaultSettings.Port
	}
	if s.Host == "-1" {
		s.Host = DefaultSettings.Host
	}
	if s.HueRate == -1.0 {
		s.HueRate = DefaultSettings.HueRate
	}
	if s.AmpScalar == -1.0 {
		s.AmpScalar = DefaultSettings.AmpScalar
	}
	if s.FilterRange == -1.0 {
		s.FilterRange = DefaultSettings.FilterRange
	}
	if s.FilterMode == -1.0 {
		s.FilterMode = DefaultSettings.FilterMode
	}
	if s.Decay == -1.0 {
		s.Decay = DefaultSettings.Decay
	}

	return nil
}

func (s *Settings) getSettingsDefaultPath() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	vbzConfigPath := path.Join(configDir, "vbz", "config")

	f, err := os.Open(vbzConfigPath)

	// will asume that the error is "no such file"
	if err != nil {
		return nil
	}

	return s.parseConfigFile(f)
}

func (s *Settings) getSettings(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.New("error while trying to open specified config file: " + err.Error())
	}

	return s.parseConfigFile(f)
}

func (s *Settings) parseConfigFile(file *os.File) error {
	buf, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	conts := string(buf)

	// ignore whitespaces
	conts = strings.ReplaceAll(conts, " ", "")
	lines := strings.Split(conts, "\n")

	var (
		funcType interface{}
		ok       bool
	)

	for i, l := range lines {
		// ignore empty lines
		if len(l) == 0 {
			continue
		}
		// skip comments
		if l[0] == '#' {
			continue
		}
		keyVal := strings.Split(l, "=")

		if len(keyVal) != 2 {
			return errors.New(fmt.Sprintf("config error near line: %d", i+1))
		}

		key := keyVal[0]
		val := keyVal[1]

		if funcType, ok, key = getFuncType(key); !ok {
			return errors.New(fmt.Sprintf("no such key on line: %d: %s", i+1, keyVal[0]))
		}

		switch f := funcType.(type) {
		case setIntConfigType:
			err = f(getFieldPointer(s, key).(*int), val, i+1)
		case setStringConfigType:
			f(getFieldPointer(s, key).(*string), val)
		case setBoolConfigType:
			err = f(getFieldPointer(s, key).(*bool), val, i+1)
		case setFloatConfigType:
			err = f(getFieldPointer(s, key).(*float64), val, i+1)
		case setFilterModeConfigType:
			err = f(getFieldPointer(s, key).(*ft.FilterType), val, i+1)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
