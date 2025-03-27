package settings

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

// in config file keys are case insensitive
type Settings struct {
	Port      int
	Host      string
	DeviceIdx int
}

var DefaultSettings Settings = Settings{
	Port:      6742,
	Host:      "localhost",
	DeviceIdx: 0,
}

func GetSettingsDefaultPath() Settings {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return DefaultSettings
	}

	vbzConfigPath := path.Join(configDir, "vbz", "config")

	f, err := os.Open(vbzConfigPath)

	// will asume that the error is "no such file"
	if err != nil {
		return DefaultSettings
	}

	return parseConfigFile(f)
}

func GetSettings(path string) Settings {
	f, err := os.Open(path)
	// will asume that the error is "no such file"
	if err != nil {
		panic("error while trying to open specified config file: " + err.Error())
	}

	return parseConfigFile(f)
}

func parseInt(s string, lineNum int) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("error while parsing int at line: %d", lineNum))
	}
	return num
}

func parseConfigFile(file *os.File) Settings {
	s := DefaultSettings

	buf, err := io.ReadAll(file)
	if err != nil {
		return DefaultSettings
	}

	conts := string(buf)

	// ignore whitespaces
	conts = strings.ReplaceAll(conts, " ", "")
	lines := strings.Split(conts, "\n")

	for i, l := range lines {
		if len(l) == 0 {
			continue
		}
		keyVal := strings.Split(l, "=")

		if len(keyVal) != 2 {
			panic(fmt.Sprintf("config error near line: %d", i+1))
		}

		key := strings.ToLower(keyVal[0])
		val := keyVal[1]

		switch key {
		case "port":
			s.Port = parseInt(val, i+1)
		case "host":
			s.Host = val
		case "deviceidx":
			s.DeviceIdx = parseInt(val, i+1)
		default:
			panic(fmt.Sprintf("no such key on line: %d", i+1))
		}
	}

	return s
}
