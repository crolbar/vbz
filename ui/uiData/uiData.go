package uiData

import (
	"time"
	"vbz/bpm"
	"vbz/fft"
	"vbz/hues"
	"vbz/settings"
)

type UiData struct {
	Fft  *fft.FFT
	Sets *settings.Settings
	Hues *hues.Hues
	Bpm  *bpm.BPM

	TickCount    uint
	LastTickTime time.Time
	FPS          int
}
