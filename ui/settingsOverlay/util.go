package settingsOverlay

import (
	"vbz/audioCapture"

	lbb "github.com/crolbar/lipbalm/components/button"
	lbl "github.com/crolbar/lipbalm/layout"
)

func (o *SettingsOverlay) initBDevices() {
	_, devices, _ := audioCapture.GetDevices()
	for i := 0; i < o.d.Audio.NumDevices; i++ {
		o.bDevices[i] = lbb.NewButtonR(devices[i].Name(), lbl.NewRect(0, 0, 1, 1), lbb.WithBorder())
	}
}

func (o *SettingsOverlay) updateRects() {
	o.bNoLeds.Rect = o.rects[bNoLeds]
	o.bFillBins.Rect = o.rects[bFillBins]
	o.bSetBlack.Rect = o.rects[bSetBlack]
	o.tiHost.Rect = o.rects[tiHost]
	o.tiPort.Rect = o.rects[tiPort]
	o.tiAmpScalar.Rect = o.rects[tiAmpScalar]
	o.tiDecay.Rect = o.rects[tiDecay]
	o.tiFilterRange.Rect = o.rects[tiFilterRange]
	o.tiHueRate.Rect = o.rects[tiHueRate]
	o.sAmpScalar.Rect = o.rects[sAmpScalar]
	o.sHueRate.Rect = o.rects[sHueRate]

	for i := FilterModeButtonsOffset; i < DeviceButtonsOffset; i++ {
		o.bFilterModes[i-FilterModeButtonsOffset].Rect = o.rects[i]
	}

	for i := DeviceButtonsOffset; i < len(o.rects); i++ {
		o.bDevices[i-DeviceButtonsOffset].Rect = o.rects[i]
	}
}
