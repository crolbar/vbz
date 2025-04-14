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
		o.comps = append(o.comps,
			lbb.Init(devices[i].Name(),
				lbb.WithBorder(),
				lbb.WithTrigger(o.handleBDevices, i),
			))
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
