package settingsOverlay

import (
	lbc "github.com/crolbar/lipbalm/components"
	lbti "github.com/crolbar/lipbalm/components/textInput"
)

func (o *SettingsOverlay) wrapFS(fc lbc.Component) func(any) error {
	return func(any) error { o.handleFocusSwitch(fc); return nil }
}

func (o *SettingsOverlay) handleFocusSwitchKbPrev() {
	o.focusedComponentKb = min(o.compsLen-1, o.focusedComponentKb+1)
	o.handleFocusSwitch(o.comps[o.focusedComponentKb])
}

func (o *SettingsOverlay) handleFocusSwitchKbNext() {
	o.focusedComponentKb = max(0, o.focusedComponentKb-1)
	o.handleFocusSwitch(o.comps[o.focusedComponentKb])
}

func (o *SettingsOverlay) handleFocusSwitch(fc lbc.Component) {
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
		return o.handleTiAmpScalar(nil)
	case tiHueRateTitle:
		return o.handleTiHueRate(nil)
	case tiHostTitle:
		return o.handleTiHost(nil)
	case tiPortTitle:
		return o.handleTiPort(nil)
	case tiFilterRangeTitle:
		return o.handleTiFilterRange(nil)
	case tiDecayTitle:
		return o.handleTiDecay(nil)
	}
	return nil
}
