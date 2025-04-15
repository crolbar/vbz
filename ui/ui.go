package ui

import (
	"time"
	"vbz/audioCapture"
	"vbz/bpm"
	"vbz/fft"
	"vbz/hues"
	"vbz/led"
	"vbz/settings"
	"vbz/ui/settingsOverlay"
	"vbz/ui/tab"
	"vbz/ui/tab/bins"
	"vbz/ui/tab/circle"
	"vbz/ui/tab/master"
	"vbz/ui/uiData"

	tea "github.com/charmbracelet/bubbletea"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
)

type Ui struct {
	d uiData.UiData

	Width  int
	Height int

	fb       lbfb.FrameBuffer
	cachedFb string

	tabs   [tab.Last__]tab.Tab
	selTab tab.TabType

	SettingsOverlay *settingsOverlay.SettingsOverlay
	showOverlay     bool
}

func InitUi(
	fft *fft.FFT,
	sets *settings.Settings,
	hues *hues.Hues,
	bpm *bpm.BPM,
	led *led.LED,
	audio *audioCapture.AudioCapture,
) Ui {
	u := Ui{
		d: uiData.UiData{
			Fft:   fft,
			Sets:  sets,
			Hues:  hues,
			Bpm:   bpm,
			Led:   led,
			Audio: audio,

			TickCount:    0,
			LastTickTime: time.Time{},
			FPS:          0,
		},

		fb: lbfb.NewFrameBuffer(0, 0),

		selTab:      tab.Master,
		showOverlay: false,
	}

	u.SettingsOverlay = settingsOverlay.Init(u.d)
	u.tabs[tab.Bins] = bins.Init(u.d)
	u.tabs[tab.Circle] = circle.Init(u.d)
	u.tabs[tab.Master] = master.Init(u.d, u.tabs[tab.Bins], u.tabs[tab.Circle])

	return u
}

func (ui *Ui) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !ui.showOverlay {
			switch msg.String() {
			case "1":
				ui.selTab = tab.Master
			case "2":
				ui.selTab = tab.Bins
			case "3":
				ui.selTab = tab.Circle
			case "q":
				return tea.Quit
			}
		}

		switch msg.String() {
		case "`", "f1":
			ui.showOverlay = !ui.showOverlay
		}
	}

	if ui.showOverlay {
		ui.SettingsOverlay.Update(msg)
	}

	ui.SelTab().Update(msg)

	if ui.PreView() {
		ui.d.TickCount++
		now := time.Now()
		frameTime := now.Sub(ui.d.LastTickTime)
		ui.d.LastTickTime = now
		ui.d.FPS = int(time.Second / frameTime)
	}

	return nil
}

func (ui *Ui) SetSize(msg tea.WindowSizeMsg) {
	ui.Width = msg.Width
	ui.Height = msg.Height
	ui.fb.Resize(ui.Width, ui.Height)

	ui.SettingsOverlay.Resize(msg)
	ui.SelTab().Resize(msg)
}

func (ui *Ui) SelTab() tab.Tab {
	return ui.tabs[ui.selTab]
}

func (ui *Ui) PreView() (c bool) {
	if ui.Width == 0 || ui.Height == 0 {
		ui.cachedFb = ""
		return
	}
	if len(ui.d.Fft.Bins) == 0 {
		ui.cachedFb = "no fft bins"
		return
	}

	if ui.d.Fft.PeakAmp != 0 && !ui.d.Fft.IsBinsUpdated() {
		return
	}

	copy(ui.d.Fft.PrevBins[:], ui.d.Fft.Bins)

	ui.fb.Clear()

	ui.SelTab().Render(&ui.fb)

	if ui.showOverlay {
		ui.SettingsOverlay.Render(&ui.fb)
	}

	if ui.d.Sets.Debug {
		ui.renderDebug()
	}

	ui.cachedFb = ui.fb.View()
	return true
}

func (ui Ui) View() string {
	return ui.cachedFb
}
