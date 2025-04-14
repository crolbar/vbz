package settingsOverlay

import (
	lbc "github.com/crolbar/lipbalm/components"
	lbb "github.com/crolbar/lipbalm/components/button"
	lbs "github.com/crolbar/lipbalm/components/slider"
	lbti "github.com/crolbar/lipbalm/components/textInput"
	"vbz/audioCapture"
)

func (o *SettingsOverlay) setErrorText(err error) {
	if err != nil {
		o.errorText.SetText(err.Error())
	} else {
		o.errorText.Text.Reset()
	}
}

func (o *SettingsOverlay) initBDevices() {
	_, devices, _ := audioCapture.GetDevices()
	for i := 0; i < o.d.Audio.NumDevices; i++ {
		name := devices[i].Name()
		o.deviceNames = append(o.deviceNames, name)

		o.comps = append(o.comps,
			lbb.Init(name,
				lbb.WithBorder(),
				lbb.WithTrigger(o.handleBDevices, i),
			))
	}
}

func (o *SettingsOverlay) shortenBDevicesNames(width int) {
	for i := 0; i < o.d.Audio.NumDevices; i++ {
		name := o.deviceNames[i]
		castAsButton(o.comps[DeviceButtonsOffset+i]).Border.Text = name[:min(width, len(name)-1)]
	}
}

// assuming that c is button
func castAsButton(c lbc.Component) *lbb.Button {
	return c.(*lbb.Button)
}

func castAsSlider(c lbc.Component) *lbs.Slider {
	return c.(*lbs.Slider)
}

func castAsTi(c lbc.Component) *lbti.TextInput {
	return c.(*lbti.TextInput)
}
