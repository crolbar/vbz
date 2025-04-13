package main

import (
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
	lb "github.com/crolbar/lipbalm"
	"github.com/gen2brain/malgo"
)

type VBZ struct {
	p     *tea.Program
	led   *led.LED
	audio *audioCapture.AudioCapture
	bpm   *bpm.BPM
	fft   *fft.FFT
	hues  *hues.Hues

	ui ui.Ui

	settings *settings.Settings

	shouldNotEnterTui bool
}

const SampleRate float64 = 10000

func main() {
	vbz, err := initVBZ()
	if err != nil {
		fmt.Println(lb.SetColor(lb.Color(1), err.Error()))
		return
	}
	if vbz.shouldNotEnterTui {
		return
	}

	vbz.audio.StartDev()

	*vbz.p = *tea.NewProgram(vbz, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := vbz.p.Run(); err != nil {
		fmt.Println(err)
		return
	}

	vbz.audio.Dev.Uninit()
	if vbz.led.Conn != nil {
		vbz.led.Conn.Close()
	}
}

func (v VBZ) Init() tea.Cmd {
	return nil
}

func (v VBZ) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return v, tea.Quit
		}
	case tea.WindowSizeMsg:
		v.ui.SetSize(msg)
	case Refresh:
	}

	cmd = v.ui.Update(msg)

	return v, cmd
}

func (v VBZ) View() string {
	return v.ui.View()
}

func (v *VBZ) onData() malgo.DataProc {
	return func(pOutputSample, pInputSamples []byte, frameCount uint32) {
		samples, _ := byteToU8(pInputSamples)
		v.fft.UpdateFFT(samples)
		v.fft.UpdatePeakLowAmp()
		v.bpm.UpdateBPM(samples)
		v.hues.UpdateHues(v.settings.HueRate, v.bpm.Bpm)

		if !v.settings.NoLeds && v.led.Conn != nil {
			v.led.SetVibe(v.hues, v.fft.PeakLowAmp)
		}

		v.triggerRefresh()
		time.Sleep(1 * time.Millisecond)
	}
}
