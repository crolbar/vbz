package settings

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"vbz/fft"
)

// in config file keys are case insensitive
type Settings struct {
	Port      int
	Host      string
	DeviceIdx int
	FftPtr    *fft.FFT
}

var DefaultSettings Settings = Settings{
	Port:      6742,
	Host:      "localhost",
	DeviceIdx: 0,
}

func (s *Settings) InitSettings(configPath string) error {
	var err error
	if configPath == "" {
		err = s.getSettingsDefaultPath()
	} else {
		err = s.getSettings(configPath)
	}
	if err != nil {
		return err
	}

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

func parseInt(s string, lineNum int) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("error while parsing int at line: %d", lineNum))
	}
	return num, nil
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

		key := strings.ToLower(keyVal[0])
		val := keyVal[1]

		switch key {
		case "ampscalar", "amp-scalar":
			if s.FftPtr.AmpScalar != fft.DefaultFFT.AmpScalar {
				break
			}

			num, err := parseInt(val, i+1)
			if err != nil {
				return err
			}
			s.FftPtr.AmpScalar = num
		case "filtermode", "filter-mode":
			if s.FftPtr.FilterMode != fft.DefaultFFT.FilterMode {
				break
			}

			switch strings.ToLower(val) {
			case "block":
				s.FftPtr.FilterMode = fft.Block
			case "box", "boxfilter":
				s.FftPtr.FilterMode = fft.BoxFilter
			case "dbox", "doublebox", "doubleboxfilter":
				s.FftPtr.FilterMode = fft.DoubleBoxFilter
			case "none":
				s.FftPtr.FilterMode = fft.None
			default:
				return errors.New("no such filter specifed in config file: " + val)
			}
		case "filterrange", "filter-range":
			if s.FftPtr.FilterRange != fft.DefaultFFT.FilterRange {
				break
			}

			num, err := parseInt(val, i+1)
			if err != nil {
				return err
			}
			s.FftPtr.FilterRange = num
		case "decay":
			if s.FftPtr.Decay != fft.DefaultFFT.Decay {
				break
			}

			num, err := parseInt(val, i+1)
			if err != nil {
				return err
			}
			s.FftPtr.Decay = num
		case "port":
			if s.Port != -1 {
				break
			}

			num, err := parseInt(val, i+1)
			if err != nil {
				return err
			}
			s.Port = num
		case "host":
			if s.Host != "-1" {
				break
			}

			s.Host = val
		case "deviceidx":
			if s.DeviceIdx != -1 {
				break
			}

			num, err := parseInt(val, i+1)
			if err != nil {
				return err
			}
			s.DeviceIdx = num
		default:
			return errors.New(fmt.Sprintf("no such key on line: %d", i+1))

		}
	}

	return nil
}
