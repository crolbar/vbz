package settingsOverlay

import (
	"fmt"
	"vbz/ui/uiData"

	tea "github.com/charmbracelet/bubbletea"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
	lbl "github.com/crolbar/lipbalm/layout"
	ft "vbz/fft/filter_types"

	lb "github.com/crolbar/lipbalm"
	lbb "github.com/crolbar/lipbalm/components/button"
	lbht "github.com/crolbar/lipbalm/components/hitTesting"
	lbs "github.com/crolbar/lipbalm/components/slider"
	lbti "github.com/crolbar/lipbalm/components/textInput"
)

type selectedRect int

const (
	bNoLeds selectedRect = iota
	bFillBins
	bSetBlack
	tiHost
	tiPort
	tiAmpScalar
	tiDecay
	tiFilterRange
	tiHueRate
	sAmpScalar
	sHueRate
	Last__
)

type SettingsOverlay struct {
	d uiData.UiData

	bFilterModes [ft.Last__]lbb.Button
	bDevices     []lbb.Button
	bNoLeds      lbb.Button
	bFillBins    lbb.Button
	bSetBlack    lbb.Button

	tiHost        lbti.TextInput
	tiPort        lbti.TextInput
	tiAmpScalar   lbti.TextInput
	tiDecay       lbti.TextInput
	tiFilterRange lbti.TextInput
	tiHueRate     lbti.TextInput

	sAmpScalar lbs.Slider
	sHueRate   lbs.Slider

	focusedComponent focusedComponent

	ht lbht.HitTesting

	// from 0 to Last__-1 are single rects
	// from Last__ to Last__+ft.Last__-1 are filter mode button rects
	// from Last__+ft.Last__ to Last__+ft.Last__+d.Audio.NumDevices-1 are device buttons
	rects       []lbl.Rect
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
		rectsSize = int(Last__) + int(ft.Last__) + d.Audio.NumDevices

		o = &SettingsOverlay{
			d:        d,
			rects:    make([]lbl.Rect, rectsSize),
			bDevices: make([]lbb.Button, d.Audio.NumDevices),
			ht:       lbht.InitHT(rectsSize),
			bNoLeds: lbb.NewButtonR("NoLeds", lbl.NewRect(0, 0, 1, 1),
				lbb.WithBorder(),
				lbb.WithInitState(d.Sets.NoLeds),
			),
			bFillBins: lbb.NewButtonR("FillBins", lbl.NewRect(0, 0, 1, 1),
				lbb.WithBorder(),
				lbb.WithInitState(d.Sets.FillBins),
			),
			bSetBlack: lbb.NewButtonR("SetBlack", lbl.NewRect(0, 0, 1, 1), lbb.WithBorder()),
			bFilterModes: [int(ft.Last__)]lbb.Button{
				lbb.NewButtonR("None", lbl.Rect{}, lbb.WithBorder()),
				lbb.NewButtonR("Block", lbl.Rect{}, lbb.WithBorder()),
				lbb.NewButtonR("BoxFilter", lbl.Rect{}, lbb.WithBorder()),
				lbb.NewButtonR("DoubleBoxFilter", lbl.Rect{}, lbb.WithBorder()),
			},
			sAmpScalar: lbs.NewSliderR("AmpScalar", lbl.Rect{},
				lbs.WithBorder(),
				lbs.WithInitProgress(uint8(255.0*min(float64(d.Sets.AmpScalar)/float64(MaxAmpScalar), 1.0))),
			),
			sHueRate: lbs.NewSliderR("HueRate", lbl.Rect{},
				lbs.WithBorder(),
				lbs.WithInitProgress(uint8((d.Sets.HueRate*10.0)*255.0)),
			),
			tiAmpScalar: lbti.NewTextInputR(tiAmpScalarTitle, lbl.Rect{},
				lbti.WithBorder(),
				lbti.WithInitText(fmt.Sprintf("%d", d.Sets.AmpScalar)),
			),
			tiHueRate: lbti.NewTextInputR(tiHueRateTitle, lbl.Rect{},
				lbti.WithBorder(),
				lbti.WithInitText(fmt.Sprintf("%.4f", d.Sets.HueRate)),
			),
			tiHost: lbti.NewTextInputR(tiHostTitle, lbl.Rect{},
				lbti.WithBorder(),
				lbti.WithInitText(d.Sets.Host),
			),
			tiPort: lbti.NewTextInputR(tiPortTitle, lbl.Rect{},
				lbti.WithBorder(),
				lbti.WithInitText(fmt.Sprintf("%d", d.Sets.Port)),
			),
			tiDecay: lbti.NewTextInputR(tiDecayTitle, lbl.Rect{},
				lbti.WithBorder(),
				lbti.WithInitText(fmt.Sprintf("%d", d.Sets.Decay)),
			),
			tiFilterRange: lbti.NewTextInputR(tiFilterRangeTitle, lbl.Rect{},
				lbti.WithBorder(),
				lbti.WithInitText(fmt.Sprintf("%d", d.Sets.FilterRange)),
			),
		}
	)

	o.initBDevices()
	o.setTriggers()

	o.bFilterModes[int(o.d.Sets.FilterMode)].Press()
	o.bDevices[o.d.Sets.DeviceIdx].Press()
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

		filterRowHeight = 3 * int(ft.Last__)

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
			).Split(cols[1])
	)

	o.overlayRect = rect

	o.minHeight = filterRowHeight + 3 + 3 + 3 + 6
	o.minWidth = 3 + 3 + 3 + 3

	// l 1 row
	o.rects[sAmpScalar] = lRows[0]

	// l 2 row
	o.rects[sHueRate] = lSecondRow[0]
	o.rects[tiAmpScalar] = lSecondRow[2]
	o.rects[tiHueRate] = lSecondRow[3]

	// l 3 row
	o.rects[tiHost] = lThirdRow[0]
	o.rects[tiPort] = lThirdRow[2]
	o.rects[tiFilterRange] = lThirdRow[3]
	o.rects[tiDecay] = lThirdRow[4]

	// l 4 row
	o.rects[bNoLeds] = lFourthRow[0]
	o.rects[bFillBins] = lFourthRow[1]
	o.rects[bSetBlack] = lFourthRow[4]

	// filter modes
	for i := FilterModeButtonsOffset; i < DeviceButtonsOffset; i++ {
		var r = lRows[4]
		r.Height = 3
		r.Y += uint16(i-FilterModeButtonsOffset) * 3
		o.rects[i] = r
	}

	// devices
	for i := DeviceButtonsOffset; i < len(o.rects); i++ {
		var r = rRows[0]
		r.Height = 3
		r.Y += uint16(i-DeviceButtonsOffset) * 3
		o.rects[i] = r
	}

	o.updateRects()
}

func (o *SettingsOverlay) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if o.focusedComponent == nil {
			switch msg.String() {
			case "r":
				o.d.Led.SetAllLEDsToColor(255, 0, 0)
			case "g":
				o.d.Led.SetAllLEDsToColor(0, 255, 0)
			case "b":
				o.d.Led.SetAllLEDsToColor(0, 0, 255)
			case "B":
				o.d.Led.SetAllLEDsToColor(0, 0, 0)
			}
		}

		switch msg.String() {
		case "esc", "enter":
			o.deFocusComponent()
		}

		o.tiAmpScalar.Update(msg.String())
		o.tiHueRate.Update(msg.String())
		o.tiHost.Update(msg.String())
		o.tiPort.Update(msg.String())
		o.tiFilterRange.Update(msg.String())
		o.tiDecay.Update(msg.String())
	case tea.MouseMsg:
		if msg.String() == "left release" {
			return
		}

		o.sAmpScalar.UpdateMouseClick(msg.String(), msg.X, msg.Y, o.sAmpScalar.Rect)
		o.sHueRate.UpdateMouseClick(msg.String(), msg.X, msg.Y, o.sHueRate.Rect)
		o.ht.CheckHit(msg.X, msg.Y, o.rects[:])
	}
}

func (o SettingsOverlay) Render(fb *lbfb.FrameBuffer) {
	if o.minHeight > int(o.overlayRect.Height) || o.minWidth > int(o.overlayRect.Width) {
		o.RenderMinDimentions(fb)
		return
	}

	fb.RenderString(o.sAmpScalar.View(), o.sAmpScalar.GetRect())
	fb.RenderString(o.sHueRate.View(), o.sHueRate.GetRect())

	fb.RenderString(o.tiAmpScalar.View(), o.tiAmpScalar.GetRect())
	fb.RenderString(o.tiHueRate.View(), o.tiHueRate.GetRect())

	fb.RenderString(o.bNoLeds.View(), o.bNoLeds.GetRect())
	fb.RenderString(o.bFillBins.View(), o.bFillBins.GetRect())
	fb.RenderString(o.bSetBlack.View(), o.bSetBlack.GetRect())

	fb.RenderString(o.tiHost.View(), o.tiHost.Rect)
	fb.RenderString(o.tiPort.View(), o.tiPort.Rect)
	fb.RenderString(o.tiFilterRange.View(), o.tiFilterRange.Rect)
	fb.RenderString(o.tiDecay.View(), o.tiDecay.Rect)

	for i := 0; i < DeviceButtonsOffset-FilterModeButtonsOffset; i++ {
		fb.RenderString(o.bFilterModes[i].View(), o.bFilterModes[i].GetRect())
	}

	for i := 0; i < len(o.rects)-DeviceButtonsOffset; i++ {
		fb.RenderString(o.bDevices[i].View(), o.bDevices[i].GetRect())
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
