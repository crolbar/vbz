package ui

import (
	"time"
	"vbz/bpm"
	"vbz/fft"
	"vbz/hues"
	"vbz/settings"

	tea "github.com/charmbracelet/bubbletea"

	// lbb "github.com/crolbar/lipbalm/components/button"
	// lbht "github.com/crolbar/lipbalm/components/hitTesting"
	// lbs "github.com/crolbar/lipbalm/components/slider"
	// lbti "github.com/crolbar/lipbalm/components/textInput"
	// lbl "github.com/crolbar/lipbalm/layout"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
)

type Ui struct {
	fft  *fft.FFT
	sets *settings.Settings
	hues *hues.Hues
	bpm  *bpm.BPM

	TickCount    uint
	LastTickTime time.Time
	FPS          int

	Width  int
	Height int

	fb lbfb.FrameBuffer
}

func InitUi(
	fft *fft.FFT,
	sets *settings.Settings,
	hues *hues.Hues,
	bpm *bpm.BPM,
) Ui {
	return Ui{
		fft:  fft,
		sets: sets,
		hues: hues,
		bpm:  bpm,

		TickCount:    0,
		LastTickTime: time.Time{},
		FPS:          0,
		fb:           lbfb.NewFrameBuffer(0, 0),
	}
}

func (ui *Ui) SetSize(msg tea.WindowSizeMsg) {
	ui.Width = msg.Width
	ui.Height = msg.Height
	ui.fb.Resize(ui.Width, ui.Height)
}

func (ui Ui) View() string {
	if ui.Width == 0 || ui.Height == 0 {
		return ""
	}
	if len(ui.fft.Bins) == 0 {
		return "no fft bins"
	}

	ui.fb.Clear()

	// return ui.renderCircle()
	return ui.renderBins()
}
