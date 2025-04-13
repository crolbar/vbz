package settingsOverlay

import (
	lbti "github.com/crolbar/lipbalm/components/textInput"
	ft "vbz/fft/filter_types"
)

type focusedComponent interface {
	Focus()
	DeFocus()
	HasFocus() bool
}

func (o *SettingsOverlay) wrapFS(fc focusedComponent) func(any) error {
	return func(any) error { o.handleFocusSwitch(fc); return nil }
}

func (o *SettingsOverlay) getComponentByIdx(idx int) focusedComponent {
	// TODO FIX THIS
	switch idx {
	case int(sAmpScalar):
		return &o.sAmpScalar
	case int(sHueRate):
		return &o.sHueRate
	case int(tiAmpScalar):
		return &o.tiAmpScalar
	case int(tiHueRate):
		return &o.tiHueRate
	case int(tiHost):
		return &o.tiHost
	case int(tiPort):
		return &o.tiPort
	case int(tiFilterRange):
		return &o.tiFilterRange
	case int(tiDecay):
		return &o.tiDecay
	case int(bNoLeds):
		return &o.bFilterModes[ft.None]
	case int(bFillBins):
		return &o.bFilterModes[ft.Block]
	case int(bSetBlack):
		return &o.bFilterModes[ft.BoxFilter]
	case int(Last__) + int(ft.None):
		return &o.bFilterModes[ft.DoubleBoxFilter]
	case int(Last__) + int(ft.Block):
		return &o.bNoLeds
	case int(Last__) + int(ft.BoxFilter):
		return &o.bFillBins
	case int(Last__) + int(ft.DoubleBoxFilter):
		return &o.bSetBlack
	}

	for i := DeviceButtonsOffset; i < len(o.rects); i++ {
		if i == idx {
			return &o.bDevices[i-DeviceButtonsOffset]
		}
	}

	return nil
}

func (o *SettingsOverlay) handleFocusSwitchKbPrev() {
	o.focusedComponentKb = min(len(o.rects)-1, o.focusedComponentKb+1)
	o.handleFocusSwitch(o.getComponentByIdx(o.focusedComponentKb))
}

func (o *SettingsOverlay) handleFocusSwitchKbNext() {
	o.focusedComponentKb = max(0, o.focusedComponentKb-1)
	o.handleFocusSwitch(o.getComponentByIdx(o.focusedComponentKb))
}

func (o *SettingsOverlay) handleFocusSwitch(fc focusedComponent) {
	if o.focusedComponent == fc {
		return
	}

	if o.focusedComponent != nil {
		o.focusedComponent.DeFocus()
	}
	o.focusedComponent = fc
	o.focusedComponent.Focus()
}

func (o *SettingsOverlay) deFocusComponent() (err error) {
	if c, ok := o.focusedComponent.(*lbti.TextInput); ok {
		err = o.handleTiDefocus(c)
	}

	if o.focusedComponent != nil {
		o.focusedComponent.DeFocus()
		o.focusedComponent = nil
	}
	return
}

func (o *SettingsOverlay) handleTiDefocus(c *lbti.TextInput) error {
	// TODO: FIX THIS ? not bad but, a title is not an id
	switch c.Title {
	case tiAmpScalarTitle:
		return o.handleTiAmpScalar()
	case tiHueRateTitle:
		return o.handleTiHueRate()
	case tiHostTitle:
		return o.handleTiHost()
	case tiPortTitle:
		return o.handleTiPort()
	case tiFilterRangeTitle:
		return o.handleTiFilterRange()
	case tiDecayTitle:
		return o.handleTiDecay()
	}
	return nil
}
