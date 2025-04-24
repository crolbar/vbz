package main

import (
	"errors"
	"fmt"
	"time"
	"vbz/audioCapture"
	"vbz/bpm"
	"vbz/fft"
	"vbz/hues"
	"vbz/led"
	"vbz/settings"
	"vbz/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func initVBZ() (VBZ, error) {
	var err error

	var defaultFFT = fft.DefaultFFT
	vbz := VBZ{
		fft: &defaultFFT,
		bpm: &bpm.BPM{
			Bpm:        0,
			LastEnergy: 0,
			LastBeat:   time.Time{},
			HasBeat:    false,
		},
		settings: &settings.DefaultSettings,
		hues:     hues.InitHues(),
		p:        tea.NewProgram(VBZ{}),
	}

	// config file & defaults
	err = vbz.settings.InitSettings()
	if err != nil {
		return VBZ{}, errors.New(fmt.Sprintf("Error seting settings: %s", err.Error()))
	}
	vbz.fft.SetSettingsPtr(vbz.settings)

	// cli arguments
	vbz.shouldNotEnterTui, err = vbz.settings.ParseEarlyArgs()
	if err != nil {
		return VBZ{}, err
	}
	if vbz.shouldNotEnterTui {
		return vbz, nil
	}

	// openrgb connection
	// we need this pointer because im passing it to the ui
	// after i can just change the value behind it
	vbz.led = &led.LED{}
	if !vbz.settings.NoOpenRgb {
		vbz.led, err = led.InitLED(vbz.settings.Host, vbz.settings.Port)
		if err != nil {
			return VBZ{}, errors.New(fmt.Sprintf("Error initializing openrgb connection: %s", err.Error()))
		}
	}

	// capture device for audio samples
	audio, err := audioCapture.InitAudioCapture(
		vbz.settings.DeviceIdx,
		fft.BUFFER_SIZE,
		SampleRate,
		vbz.onData(),
	)
	if err != nil {
		return VBZ{}, errors.New(fmt.Sprintf("Error initializing audio capture dev: %s", err.Error()))
	}
	vbz.audio = &audio

	// late cli args for leds
	err = vbz.parseLateArgs()
	if err != nil {
		return VBZ{}, errors.New(fmt.Sprintf("Error parsing args: %s", err.Error()))
	}

	// ui
	vbz.ui = ui.InitUi(
		vbz.fft,
		vbz.settings,
		vbz.hues,
		vbz.bpm,
		vbz.led,
		vbz.audio,
	)

	return vbz, nil
}
