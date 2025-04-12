package settingsOverlay

import (
	"fmt"
	"time"
	ft "vbz/fft/filter_types"

	lbti "github.com/crolbar/lipbalm/components/textInput"
)

func timedAction(after time.Duration, c func()) {
	time.Sleep(after)
	c()
}

func (o *SettingsOverlay) setTriggers() {
	for i := 0; i < len(o.rects); i++ {
		if i >= DeviceButtonsOffset { // device buttons
			o.ht.SetTrigger(i, o.handleBDevices(i-DeviceButtonsOffset))
			continue
		}

		if i >= FilterModeButtonsOffset { // filtermode buttons
			o.ht.SetTrigger(i, o.handleBFilterModes(ft.FilterType(i-FilterModeButtonsOffset)))
			continue
		}

		if i < FilterModeButtonsOffset { // single rects from const vals
			switch selectedRect(i) {
			case bNoLeds:
				o.ht.SetTrigger(i, o.handleBNoLedsTrigger)
			case bFillBins:
				o.ht.SetTrigger(i, o.handleBFillBinsTrigger)
			case bSetBlack:
				o.ht.SetTrigger(i, o.handleBSetBlackTrigger)
			case tiHost:
				o.ht.SetTrigger(i, func(any) { o.handleFocusSwitch(&o.tiHost) })
			case tiPort:
				o.ht.SetTrigger(i, func(any) { o.handleFocusSwitch(&o.tiPort) })
			case tiAmpScalar:
				o.ht.SetTrigger(i, func(any) { o.handleFocusSwitch(&o.tiAmpScalar) })
			case tiDecay:
				o.ht.SetTrigger(i, func(any) { o.handleFocusSwitch(&o.tiDecay) })
			case tiFilterRange:
				o.ht.SetTrigger(i, func(any) { o.handleFocusSwitch(&o.tiFilterRange) })
			case tiHueRate:
				o.ht.SetTrigger(i, func(any) { o.handleFocusSwitch(&o.tiHueRate) })
			case sAmpScalar:
				o.ht.SetTrigger(i, o.handleSAmpScalar())
			case sHueRate:
				o.ht.SetTrigger(i, o.handleSHueRate())
			}
		}
	}
}

func (o *SettingsOverlay) handleTiDefocus(c *lbti.TextInput) {
	// TODO: FIX THIS ? not bad but, a title is not an id
	switch c.Title {
	case tiAmpScalarTitle:
		o.handleTiAmpScalar()
	case tiHueRateTitle:
		o.handleTiHueRate()
	case tiHostTitle:
		o.handleTiHost()
	case tiPortTitle:
		o.handleTiPort()
	case tiFilterRangeTitle:
		o.handleTiFilterRange()
	case tiDecayTitle:
		o.handleTiDecay()
	}
}

func (o *SettingsOverlay) handleTiHost() {
	o.d.Sets.Host = o.tiHost.GetText()
}

func (o *SettingsOverlay) handleTiPort() {
	n, _ := o.tiPort.GetTextAsInt()
	o.d.Sets.Port = n
}

func (o *SettingsOverlay) handleTiFilterRange() {
	n, _ := o.tiFilterRange.GetTextAsInt()
	o.d.Sets.FilterRange = n
}

func (o *SettingsOverlay) handleTiDecay() {
	n, _ := o.tiDecay.GetTextAsInt()
	o.d.Sets.Decay = min(n, 99)
}

func (o *SettingsOverlay) handleTiHueRate() {
	f, _ := o.tiHueRate.GetTextAsFloat()
	o.d.Sets.HueRate = f
	o.sHueRate.SetRatio(uint8(255 * float64(f) * 10))
}

func (o *SettingsOverlay) handleTiAmpScalar() {
	i, _ := o.tiAmpScalar.GetTextAsInt()
	o.d.Sets.AmpScalar = i
	o.sAmpScalar.SetRatio(uint8(255 * float64(i) / MaxAmpScalar))
}

func (o *SettingsOverlay) handleSHueRate() func(any) {
	return func(any) {
		o.handleFocusSwitch(&o.sHueRate)

		o.d.Sets.HueRate = o.sHueRate.GetRatio() / 10
		o.tiHueRate.SetText(fmt.Sprintf("%.4f", o.d.Sets.HueRate))
	}
}

func (o *SettingsOverlay) handleSAmpScalar() func(any) {
	return func(any) {
		o.handleFocusSwitch(&o.sAmpScalar)

		o.d.Sets.AmpScalar = int(MaxAmpScalar * o.sAmpScalar.GetRatio())
		o.tiAmpScalar.SetText(fmt.Sprintf("%d", o.d.Sets.AmpScalar))
	}
}

func (o *SettingsOverlay) handleBDevices(devIdx int) func(any) {
	return func(any) {
		o.bDevices[o.d.Sets.DeviceIdx].Depress()
		o.d.Sets.DeviceIdx = devIdx
		o.bDevices[o.d.Sets.DeviceIdx].Press()

		err := o.d.Audio.ReinitDevice(o.d.Sets.DeviceIdx)

		// TODO FIX THIS
		if err != nil {
			panic(err)
		}

		o.d.Audio.StartDev()
	}
}

func (o *SettingsOverlay) handleBFilterModes(filterT ft.FilterType) func(any) {
	return func(any) {
		o.bFilterModes[int(o.d.Sets.FilterMode)].Depress()
		o.d.Sets.FilterMode = filterT
		o.bFilterModes[int(o.d.Sets.FilterMode)].Press()
	}
}

func (o *SettingsOverlay) handleBNoLedsTrigger(any) {
	o.d.Sets.NoLeds = !o.d.Sets.NoLeds
	o.bNoLeds.Pressed = o.d.Sets.NoLeds
}

func (o *SettingsOverlay) handleBFillBinsTrigger(any) {
	o.d.Sets.FillBins = !o.d.Sets.FillBins
	o.bFillBins.Pressed = o.d.Sets.FillBins
}

func (o *SettingsOverlay) handleBSetBlackTrigger(any) {
	o.d.Led.SetAllLEDsToColor(0, 0, 0)
	o.bSetBlack.Press()
	go timedAction(time.Millisecond*200, func() {
		o.bSetBlack.Depress()
	})
}
