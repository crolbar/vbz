package uiData

import (
	"time"
	"vbz/audioCapture"
	"vbz/bpm"
	"vbz/fft"
	"vbz/hues"
	"vbz/led"
	"vbz/settings"
)

type UiData struct {
	Fft   *fft.FFT
	Sets  *settings.Settings
	Hues  *hues.Hues
	Bpm   *bpm.BPM
	Led   *led.LED
	Audio *audioCapture.AudioCapture

	TickCount    uint
	LastTickTime time.Time
	FPS          int
}
