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

func (o *SettingsOverlay) setTriggers() {
	for i := 0; i < o.compsLen; i++ {
		if i > int(tiDecay) {
			o.ht.SetTriggerFromComponent(i, o.comps[i])
			continue
		}

		// only switch focus
		switch selectedRect(i) {
		case tiHost:
			o.ht.SetTrigger(i, o.wrapFS(o.comps[tiHost]))
		case tiPort:
			o.ht.SetTrigger(i, o.wrapFS(o.comps[tiPort]))
		case tiAmpScalar:
			o.ht.SetTrigger(i, o.wrapFS(o.comps[tiAmpScalar]))
		case tiDecay:
			o.ht.SetTrigger(i, o.wrapFS(o.comps[tiDecay]))
		case tiFilterRange:
			o.ht.SetTrigger(i, o.wrapFS(o.comps[tiFilterRange]))
		case tiHueRate:
			o.ht.SetTrigger(i, o.wrapFS(o.comps[tiHueRate]))
		case sAmpScalar:
			o.ht.SetTrigger(i, o.wrapFS(o.comps[sAmpScalar]))
		case sHueRate:
			o.ht.SetTrigger(i, o.wrapFS(o.comps[sHueRate]))
		}
	}
}

func (o *SettingsOverlay) handleTiHost(any) error {
	o.d.Sets.Host = castAsTi(o.comps[tiHost]).GetText()
	newLed, err := led.InitLED(o.d.Sets.Host, o.d.Sets.Port)
	if err != nil {
		return err
	}

	*o.d.Led = *newLed
	return nil
}

func (o *SettingsOverlay) handleTiPort(any) error {
	n, err := castAsTi(o.comps[tiPort]).GetTextAsInt()
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

func (o *SettingsOverlay) handleTiFilterRange(any) error {
	n, err := castAsTi(o.comps[tiFilterRange]).GetTextAsInt()
	if err != nil {
		return err
	}
	o.d.Sets.FilterRange = n
	return nil
}

func (o *SettingsOverlay) handleTiDecay(any) error {
	n, err := castAsTi(o.comps[tiDecay]).GetTextAsInt()
	if err != nil {
		return err
	}
	o.d.Sets.Decay = min(n, 99)
	return nil
}

func (o *SettingsOverlay) handleTiHueRate(any) error {
	ti := castAsTi(o.comps[tiHueRate])
	s := castAsSlider(o.comps[sHueRate])

	f, err := ti.GetTextAsFloat()
	if err != nil {
		return err
	}
	o.d.Sets.HueRate = f
	s.SetRatio(float64(f) * 10)
	return nil
}

func (o *SettingsOverlay) handleTiAmpScalar(any) error {
	ti := castAsTi(o.comps[tiAmpScalar])
	s := castAsSlider(o.comps[sAmpScalar])

	i, err := ti.GetTextAsInt()
	if err != nil {
		return err
	}
	o.d.Sets.AmpScalar = i
	s.SetRatio(float64(i) / MaxAmpScalar)
	return nil
}

func (o *SettingsOverlay) handleSHueRate() {
	ti := castAsTi(o.comps[tiHueRate])
	s := castAsSlider(o.comps[sHueRate])

	o.d.Sets.HueRate = s.GetRatio() / 10
	ti.SetText(fmt.Sprintf("%.4f", o.d.Sets.HueRate))
}

func (o *SettingsOverlay) handleSAmpScalar() {
	ti := castAsTi(o.comps[tiAmpScalar])
	s := castAsSlider(o.comps[sAmpScalar])

	o.d.Sets.AmpScalar = int(MaxAmpScalar * s.GetRatio())
	ti.SetText(fmt.Sprintf("%d", o.d.Sets.AmpScalar))
}

func (o *SettingsOverlay) handleBDevices(a any) error {
	devIdx := a.(int)
	bPrev := castAsButton(o.comps[DeviceButtonsOffset+o.d.Sets.DeviceIdx])
	bCurr := castAsButton(o.comps[DeviceButtonsOffset+devIdx])

	bPrev.Depress()
	o.d.Sets.DeviceIdx = devIdx
	bCurr.Press()

	err := o.d.Audio.ReinitDevice(o.d.Sets.DeviceIdx)

	if err != nil {
		return err
	}

	o.d.Audio.StartDev()
	return nil
}

func (o *SettingsOverlay) handleBFilterModes(a any) error {
	filterT, ok := a.(ft.FilterType)
	if !ok {
		panic(fmt.Sprintf("not ok, v: %V", a))
	}
	bPrev := castAsButton(o.comps[FilterModeButtonsOffset+int(o.d.Sets.FilterMode)])
	bCurr := castAsButton(o.comps[FilterModeButtonsOffset+int(filterT)])

	bPrev.Depress()
	o.d.Sets.FilterMode = filterT
	bCurr.Press()
	return nil
}

func (o *SettingsOverlay) handleBNoLedsTrigger() {
	o.d.Sets.NoLeds = !o.d.Sets.NoLeds
	castAsButton(o.comps[bNoLeds]).Pressed = o.d.Sets.NoLeds
}

func (o *SettingsOverlay) handleBFillBinsTrigger() {
	o.d.Sets.FillBins = !o.d.Sets.FillBins
	castAsButton(o.comps[bFillBins]).Pressed = o.d.Sets.FillBins
}

func (o *SettingsOverlay) handleBSetBlackTrigger(any) error {
	err := o.d.Led.SetAllLEDsToColor(0, 0, 0)
	if err != nil {
		return err
	}
	b := castAsButton(o.comps[bSetBlack])

	b.Press()
	go timedAction(time.Millisecond*200, func() {
		b.Depress()
	})
	return nil
}
