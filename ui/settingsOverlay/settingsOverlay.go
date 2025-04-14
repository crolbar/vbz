package settingsOverlay

import (
	"fmt"
	"vbz/ui/uiData"

	tea "github.com/charmbracelet/bubbletea"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
	lbl "github.com/crolbar/lipbalm/layout"
	ft "vbz/fft/filter_types"

	lb "github.com/crolbar/lipbalm"
	lbc "github.com/crolbar/lipbalm/components"
	lbb "github.com/crolbar/lipbalm/components/button"
	lbht "github.com/crolbar/lipbalm/components/hitTesting"
	lbs "github.com/crolbar/lipbalm/components/slider"
	lbti "github.com/crolbar/lipbalm/components/textInput"
)

type selectedRect int

const (
	sAmpScalar selectedRect = iota
	sHueRate
	tiAmpScalar
	tiHueRate
	tiHost
	tiPort
	tiFilterRange
	tiDecay
	bNoLeds
	bFillBins
	bSetBlack
	Last__
)

type SettingsOverlay struct {
	d uiData.UiData

	deviceNames []string

	comps    []lbc.Component
	compsLen int

	errorText *lbti.TextInput

	focusedComponent lbc.Component
	// TODO something better ?
	focusedComponentKb int

	ht lbht.HitTesting

	// from 0 to Last__-1 are single rects
	// from Last__ to Last__+ft.Last__-1 are filter mode button rects
	// from Last__+ft.Last__ to Last__+ft.Last__+d.Audio.NumDevices-1 are device buttons
	// rects       []lbl.Rect
	overlayRect lbl.Rect

	minHeight int
	minWidth  int
}

const (
	FilterModeButtonsOffset int = int(Last__)
	DeviceButtonsOffset     int = int(Last__) + int(ft.Last__)

	MaxAmpScalar float64 = 10000.0

	tiAmpScalarTitle   = "AmpScalar"
	tiHueRateTitle     = "HueRate"
	tiHostTitle        = "Host"
	tiPortTitle        = "Port"
	tiDecayTitle       = "Decay"
	tiFilterRangeTitle = "FilterRange"
)

func Init(d uiData.UiData) *SettingsOverlay {
	var (
		componentsSize = int(Last__) + int(ft.Last__) + d.Audio.NumDevices

		o = &SettingsOverlay{
			d:         d,
			ht:        lbht.InitHT(componentsSize),
			errorText: lbti.Init("", lbti.WithTextColor(lb.Color(1))),

			compsLen: componentsSize,
		}
	)

	o.comps = []lbc.Component{
		sAmpScalar: lbs.Init("AmpScalar",
			lbs.WithBorder(),
			lbs.WithInitProgress(uint8(255.0*min(float64(d.Sets.AmpScalar)/float64(MaxAmpScalar), 1.0))),
			lbs.WithTrigger(wrapTrigger(o.handleSAmpScalar)),
		),
		sHueRate: lbs.Init("HueRate",
			lbs.WithBorder(),
			lbs.WithInitProgress(uint8((d.Sets.HueRate*10.0)*255.0)),
			lbs.WithTrigger(wrapTrigger(o.handleSHueRate)),
		),
		tiAmpScalar: lbti.Init(tiAmpScalarTitle,
			lbti.WithBorder(),
			lbti.WithInitText(fmt.Sprintf("%d", d.Sets.AmpScalar)),
			lbti.WithTrigger(o.handleTiAmpScalar),
		),
		tiHueRate: lbti.Init(tiHueRateTitle,
			lbti.WithBorder(),
			lbti.WithInitText(fmt.Sprintf("%.4f", d.Sets.HueRate)),
			lbti.WithTrigger(o.handleTiHueRate),
		),

		// no trigger because we don't want reconnect try on each keypress
		tiHost: lbti.Init(tiHostTitle,
			lbti.WithBorder(),
			lbti.WithInitText(d.Sets.Host),
		),
		tiPort: lbti.Init(tiPortTitle,
			lbti.WithBorder(),
			lbti.WithInitText(fmt.Sprintf("%d", d.Sets.Port)),
		),

		tiFilterRange: lbti.Init(tiFilterRangeTitle,
			lbti.WithBorder(),
			lbti.WithInitText(fmt.Sprintf("%d", d.Sets.FilterRange)),
			lbti.WithTrigger(o.handleTiFilterRange),
		),
		tiDecay: lbti.Init(tiDecayTitle,
			lbti.WithBorder(),
			lbti.WithInitText(fmt.Sprintf("%d", d.Sets.Decay)),
			lbti.WithTrigger(o.handleTiDecay),
		),
		bNoLeds: lbb.Init("NoLeds",
			lbb.WithBorder(),
			lbb.WithInitState(d.Sets.NoLeds),
			lbb.WithTrigger(wrapTrigger(o.handleBNoLedsTrigger)),
		),
		bFillBins: lbb.Init("FillBins",
			lbb.WithBorder(),
			lbb.WithInitState(d.Sets.FillBins),
			lbb.WithTrigger(wrapTrigger(o.handleBFillBinsTrigger)),
		),
		bSetBlack: lbb.Init("SetBlack",
			lbb.WithBorder(),
			lbb.WithTrigger(o.handleBSetBlackTrigger),
		),

		int(Last__) + int(ft.None):            lbb.Init("None", lbb.WithBorder(), lbb.WithTrigger(o.handleBFilterModes, ft.None)),
		int(Last__) + int(ft.Block):           lbb.Init("Block", lbb.WithBorder(), lbb.WithTrigger(o.handleBFilterModes, ft.Block)),
		int(Last__) + int(ft.BoxFilter):       lbb.Init("BoxFilter", lbb.WithBorder(), lbb.WithTrigger(o.handleBFilterModes, ft.BoxFilter)),
		int(Last__) + int(ft.DoubleBoxFilter): lbb.Init("DoubleBoxFilter", lbb.WithBorder(), lbb.WithTrigger(o.handleBFilterModes, ft.DoubleBoxFilter)),
	}

	lbs.DecreaseKeys = []string{
		"left", "h", "ctrl+b",
	}
	lbs.IncreaseKeys = []string{
		"right", "l", "ctrl+f",
	}

	o.initBDevices()
	o.setTriggers()

	castAsButton(o.comps[FilterModeButtonsOffset+int(o.d.Sets.FilterMode)]).Press()
	castAsButton(o.comps[DeviceButtonsOffset+o.d.Sets.DeviceIdx]).Press()
	return o
}

var l = lbl.DefaultLayout()

func (o *SettingsOverlay) Resize(msg tea.WindowSizeMsg) {
	var (
		w    = msg.Width
		h    = msg.Height
		xOff = uint16(0.10 * float64(w))
		yOff = uint16(0.10 * float64(h))

		rect = lbl.NewRect(xOff, yOff, uint16(w)-xOff*2, uint16(h)-yOff*2)

		cols = l.Horizontal().
			Constrains(
				lbl.NewConstrain(lbl.Percent, 50),
				lbl.NewConstrain(lbl.Percent, 50),
			).Split(rect)

		filterRowHeight = 3*int(ft.Last__) + 1

		lRows = l.Vercital().
			Constrains(
				lbl.NewConstrain(lbl.Length, 6),
				lbl.NewConstrain(lbl.Length, 3),
				lbl.NewConstrain(lbl.Length, 3),
				lbl.NewConstrain(lbl.Percent, 20),
				lbl.NewConstrain(lbl.Min, uint16(filterRowHeight)), // filters
				lbl.NewConstrain(lbl.Min, 3),                       // butons
			).Split(cols[0])

		lSecondRow = l.Horizontal().
				Constrains(
				lbl.NewConstrain(lbl.Percent, 50),
				lbl.NewConstrain(lbl.Min, 0),
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Percent, 25),
			).Split(lRows[1])

		lThirdRow = l.Horizontal().
				Constrains(
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Min, 0),
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Percent, 25),
			).Split(lRows[2])

		lFourthRow = l.Horizontal().
				Constrains(
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Min, 0),
				lbl.NewConstrain(lbl.Percent, 25),
			).Split(lRows[5])

		deviceButtonsHeight = 3 * int(o.d.Audio.NumDevices)

		rRows = l.Vercital().
			Constrains(
				lbl.NewConstrain(lbl.Length, uint16(deviceButtonsHeight)), // devices
				lbl.NewConstrain(lbl.Percent, 100),
				lbl.NewConstrain(lbl.Length, 1),
			).Split(cols[1])
	)

	o.overlayRect = rect

	o.minHeight = filterRowHeight + 3 + 3 + 3 + 6
	o.minWidth = 3 + 3 + 3 + 3

	// l 1 row
	o.comps[sAmpScalar].SetRect(lRows[0])

	// l 2 row
	o.comps[sHueRate].SetRect(lSecondRow[0])
	o.comps[tiAmpScalar].SetRect(lSecondRow[2])
	o.comps[tiHueRate].SetRect(lSecondRow[3])

	// l 3 row
	o.comps[tiHost].SetRect(lThirdRow[0])
	o.comps[tiPort].SetRect(lThirdRow[2])
	o.comps[tiFilterRange].SetRect(lThirdRow[3])
	o.comps[tiDecay].SetRect(lThirdRow[4])

	// l 4 row
	o.comps[bNoLeds].SetRect(lFourthRow[0])
	o.comps[bFillBins].SetRect(lFourthRow[1])
	o.comps[bSetBlack].SetRect(lFourthRow[4])

	// r
	o.errorText.SetRect(rRows[2])

	// filter modes
	for i := FilterModeButtonsOffset; i < DeviceButtonsOffset; i++ {
		var r = lRows[4]
		r.Height = 3
		r.Y += uint16(i-FilterModeButtonsOffset) * 3
		o.comps[i].SetRect(r)
	}

	// devices
	for i := DeviceButtonsOffset; i < o.compsLen; i++ {
		var r = rRows[0]
		r.Height = 3
		r.Y += uint16(i-DeviceButtonsOffset) * 3
		o.comps[i].SetRect(r)
		o.shortenBDevicesNames(int(r.Width) - 2)
	}
}

func (o *SettingsOverlay) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// if we are not focused on a text input
		if _, ok := o.focusedComponent.(*lbti.TextInput); !ok {
			switch msg.String() {
			case "r":
				o.d.Led.SetAllLEDsToColor(255, 0, 0)
				return
			case "g":
				o.d.Led.SetAllLEDsToColor(0, 255, 0)
				return
			case "b":
				o.d.Led.SetAllLEDsToColor(0, 0, 255)
				return
			case "B":
				o.d.Led.SetAllLEDsToColor(0, 0, 0)
				return
			case "k":
				o.handleFocusSwitchKbNext()
				return
			case "j":
				o.handleFocusSwitchKbPrev()
				return
			case "f":
				o.handleFocusSwitch(o.comps[o.focusedComponentKb])
				return
			}
		}

		switch msg.String() {
		case "alt+k", "up":
			o.handleFocusSwitchKbNext()
		case "alt+j", "down":
			o.handleFocusSwitchKbPrev()
		case "esc", "enter":
			err := o.deFocusComponent()
			o.setErrorText(err)
			return
		}

		for i := 0; i < o.compsLen; i++ {
			o.comps[i].Update(msg.String())
		}
	case tea.MouseMsg:
		if msg.String() == "left release" {
			return
		}

		castAsSlider(o.comps[sAmpScalar]).UpdateMouseClick(msg.String(), msg.X, msg.Y, o.comps[sAmpScalar].GetRect())
		castAsSlider(o.comps[sHueRate]).UpdateMouseClick(msg.String(), msg.X, msg.Y, o.comps[sHueRate].GetRect())

		err := o.ht.CheckHitOnComponents(msg.X, msg.Y, o.comps)
		o.setErrorText(err)
	}
}

func (o SettingsOverlay) Render(fb *lbfb.FrameBuffer) {
	if o.minHeight > int(o.overlayRect.Height) || o.minWidth > int(o.overlayRect.Width) {
		o.RenderMinDimentions(fb)
		return
	}

	// settings setters
	for i := 0; i < FilterModeButtonsOffset; i++ {
		fb.RenderString(o.comps[i].View(), o.comps[i].GetRect())
	}

	// filter modes
	for i := FilterModeButtonsOffset; i < DeviceButtonsOffset; i++ {
		filters := ""
		filtersRect := o.comps[FilterModeButtonsOffset].GetRect()
		filtersRect.Height *= uint16(o.d.Audio.NumDevices)
		for i := FilterModeButtonsOffset; i < DeviceButtonsOffset; i++ {
			if filters == "" {
				filters = o.comps[i].View()
			} else {
				filters = lb.JoinVertical(lb.Left, filters, o.comps[i].View())
			}
		}

		fb.RenderString(
			lb.Border(lb.NormalBorder(
				lb.WithTextBottom(
					lb.SetColor(lb.Color(135), "Avg Filter Modes"),
					lb.Center),
			), filters, true, true, false, true),
			filtersRect)
	}

	// devices
	{
		devices := ""
		devicesRect := o.comps[DeviceButtonsOffset].GetRect()
		devicesRect.Height *= uint16(o.d.Audio.NumDevices)
		for i := DeviceButtonsOffset; i < o.compsLen; i++ {
			if devices == "" {
				devices = o.comps[i].View()
			} else {
				devices = lb.JoinVertical(lb.Left, devices, o.comps[i].View())
			}
		}

		fb.RenderString(
			lb.Border(lb.NormalBorder(
				lb.WithTextBottom(
					lb.SetColor(lb.Color(135), "Capture Devices"),
					lb.Center),
			), devices, true, true, false, true),
			devicesRect)
	}

	if o.errorText.Text.Len() > 0 {
		fb.RenderString(o.errorText.View(), o.errorText.Rect)
	}
}

func (o SettingsOverlay) RenderMinDimentions(fb *lbfb.FrameBuffer) {
	s := lb.SetColor(lb.Color(1),
		fmt.Sprintf("OVERLAY WINDOW\nmin height(%d) & width(%d) reached\ncurrent h: %d, w: %d",
			o.minHeight, o.minWidth, o.overlayRect.Height, o.overlayRect.Width,
		))
	w := lb.GetWidth(s)
	fbS := fb.Size()
	fb.RenderString(
		s,
		lbl.NewRect((fbS.Width/2)-uint16(w/2), (fbS.Height/2)+2, uint16(w), 3),
		lb.Center, lb.Center)
}
