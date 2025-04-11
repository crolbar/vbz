package ui

import (
	"time"
	"vbz/bpm"
	"vbz/fft"
	"vbz/hues"
	"vbz/settings"
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

	fb lbfb.FrameBuffer

	tabs   [tab.Last__]tab.Tab
	selTab tab.TabType
}

func InitUi(
	fft *fft.FFT,
	sets *settings.Settings,
	hues *hues.Hues,
	bpm *bpm.BPM,
) Ui {
	u := Ui{
		d: uiData.UiData{
			Fft:  fft,
			Sets: sets,
			Hues: hues,
			Bpm:  bpm,

			TickCount:    0,
			LastTickTime: time.Time{},
			FPS:          0,
		},

		fb: lbfb.NewFrameBuffer(0, 0),

		selTab: tab.Master,
	}

	u.tabs[tab.Bins] = bins.Init(u.d)
	u.tabs[tab.Circle] = circle.Init(u.d)
	u.tabs[tab.Master] = master.Init(u.d, u.tabs[tab.Bins], u.tabs[tab.Circle])

	return u
}

func (ui *Ui) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			ui.selTab = tab.Master
		case "2":
			ui.selTab = tab.Bins
		case "3":
			ui.selTab = tab.Circle
		}
	}

	ui.SelTab().Update(msg)

	ui.d.TickCount++
	now := time.Now()
	frameTime := now.Sub(ui.d.LastTickTime)
	ui.d.LastTickTime = now
	ui.d.FPS = int(time.Second / frameTime)
}

func (ui *Ui) SetSize(msg tea.WindowSizeMsg) {
	ui.Width = msg.Width
	ui.Height = msg.Height
	ui.fb.Resize(ui.Width, ui.Height)

	ui.SelTab().Resize(msg)
}

func (ui *Ui) SelTab() tab.Tab {
	return ui.tabs[ui.selTab]
}

func (ui Ui) View() string {
	if ui.Width == 0 || ui.Height == 0 {
		return ""
	}
	if len(ui.d.Fft.Bins) == 0 {
		return "no fft bins"
	}

	ui.fb.Clear()

	ui.SelTab().Render(&ui.fb)

	return ui.fb.View()
}
