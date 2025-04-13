package settingsOverlay

import (
	"fmt"
	"time"
	ft "vbz/fft/filter_types"
	"vbz/led"
)

func timedAction(after time.Duration, c func()) {
	time.Sleep(after)
	c()
}

func wrapTrigger(c func()) func(any) error {
	return func(any) error {
		c()
		return nil
	}
}

func wrapTriggerE(c func() error) func(any) error {
	return func(any) error {
		return c()
	}
}

func (o *SettingsOverlay) setTriggers() {
	for i := 0; i < len(o.rects); i++ {
		if i >= DeviceButtonsOffset { // device buttons
			o.ht.SetTrigger(i, func(any) error {
				return o.handleBDevices(i - DeviceButtonsOffset)
			})
			continue
		}

		if i >= FilterModeButtonsOffset { // filtermode buttons
			o.ht.SetTrigger(i, func(any) error {
				o.handleBFilterModes(ft.FilterType(i - FilterModeButtonsOffset))
				return nil
			})
			continue
		}

		if i < FilterModeButtonsOffset { // single rects from const vals
			switch selectedRect(i) {
			case bNoLeds:
				o.ht.SetTrigger(i, wrapTrigger(o.handleBNoLedsTrigger))
			case bFillBins:
				o.ht.SetTrigger(i, wrapTrigger(o.handleBFillBinsTrigger))
			case bSetBlack:
				o.ht.SetTrigger(i, wrapTriggerE(o.handleBSetBlackTrigger))
			case tiHost:
				o.ht.SetTrigger(i, o.wrapFS(&o.tiHost))
			case tiPort:
				o.ht.SetTrigger(i, o.wrapFS(&o.tiPort))
			case tiAmpScalar:
				o.ht.SetTrigger(i, o.wrapFS(&o.tiAmpScalar))
			case tiDecay:
				o.ht.SetTrigger(i, o.wrapFS(&o.tiDecay))
			case tiFilterRange:
				o.ht.SetTrigger(i, o.wrapFS(&o.tiFilterRange))
			case tiHueRate:
				o.ht.SetTrigger(i, o.wrapFS(&o.tiHueRate))
			case sAmpScalar:
				o.ht.SetTrigger(i, wrapTrigger(o.handleSAmpScalar))
			case sHueRate:
				o.ht.SetTrigger(i, wrapTrigger(o.handleSHueRate))
			}
		}
	}
}

func (o *SettingsOverlay) handleTiHost() error {
	o.d.Sets.Host = o.tiHost.GetText()
	newLed, err := led.InitLED(o.d.Sets.Host, o.d.Sets.Port)
	if err != nil {
		return err
	}

	*o.d.Led = *newLed
	return nil
}

func (o *SettingsOverlay) handleTiPort() error {
	n, err := o.tiPort.GetTextAsInt()
	if err != nil {
		return err
	}
	o.d.Sets.Port = n

	newLed, err := led.InitLED(o.d.Sets.Host, o.d.Sets.Port)
	if err != nil {
		return err
	}

	*o.d.Led = *newLed

	return nil
}

func (o *SettingsOverlay) handleTiFilterRange() error {
	n, err := o.tiFilterRange.GetTextAsInt()
	if err != nil {
		return err
	}
	o.d.Sets.FilterRange = n
	return nil
}

func (o *SettingsOverlay) handleTiDecay() error {
	n, err := o.tiDecay.GetTextAsInt()
	if err != nil {
		return err
	}
	o.d.Sets.Decay = min(n, 99)
	return nil
}

func (o *SettingsOverlay) handleTiHueRate() error {
	f, err := o.tiHueRate.GetTextAsFloat()
	if err != nil {
		return err
	}
	o.d.Sets.HueRate = f
	o.sHueRate.SetRatio(uint8(255 * float64(f) * 10))
	return nil
}

func (o *SettingsOverlay) handleTiAmpScalar() error {
	i, err := o.tiAmpScalar.GetTextAsInt()
	if err != nil {
		return err
	}
	o.d.Sets.AmpScalar = i
	o.sAmpScalar.SetRatio(uint8(255 * float64(i) / MaxAmpScalar))
	return nil
}

func (o *SettingsOverlay) handleSHueRate() {
	o.handleFocusSwitch(&o.sHueRate)

	o.d.Sets.HueRate = o.sHueRate.GetRatio() / 10
	o.tiHueRate.SetText(fmt.Sprintf("%.4f", o.d.Sets.HueRate))
}

func (o *SettingsOverlay) handleSAmpScalar() {
	o.handleFocusSwitch(&o.sAmpScalar)

	o.d.Sets.AmpScalar = int(MaxAmpScalar * o.sAmpScalar.GetRatio())
	o.tiAmpScalar.SetText(fmt.Sprintf("%d", o.d.Sets.AmpScalar))
}

func (o *SettingsOverlay) handleBDevices(devIdx int) error {
	o.bDevices[o.d.Sets.DeviceIdx].Depress()
	o.d.Sets.DeviceIdx = devIdx
	o.bDevices[o.d.Sets.DeviceIdx].Press()

	err := o.d.Audio.ReinitDevice(o.d.Sets.DeviceIdx)

	if err != nil {
		return err
	}

	o.d.Audio.StartDev()
	return nil
}

func (o *SettingsOverlay) handleBFilterModes(filterT ft.FilterType) {
	o.bFilterModes[int(o.d.Sets.FilterMode)].Depress()
	o.d.Sets.FilterMode = filterT
	o.bFilterModes[int(o.d.Sets.FilterMode)].Press()
}

func (o *SettingsOverlay) handleBNoLedsTrigger() {
	o.d.Sets.NoLeds = !o.d.Sets.NoLeds
	o.bNoLeds.Pressed = o.d.Sets.NoLeds
}

func (o *SettingsOverlay) handleBFillBinsTrigger() {
	o.d.Sets.FillBins = !o.d.Sets.FillBins
	o.bFillBins.Pressed = o.d.Sets.FillBins
}

func (o *SettingsOverlay) handleBSetBlackTrigger() error {
	err := o.d.Led.SetAllLEDsToColor(0, 0, 0)
	if err != nil {
		return err
	}
	o.bSetBlack.Press()
	go timedAction(time.Millisecond*200, func() {
		o.bSetBlack.Depress()
	})
	return nil
}
