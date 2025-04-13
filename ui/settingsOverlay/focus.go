package settingsOverlay

import lbti "github.com/crolbar/lipbalm/components/textInput"

type focusedComponent interface {
	Focus()
	DeFocus()
	HasFocus() bool
}

func (o *SettingsOverlay) wrapFS(fc focusedComponent) func(any) error {
	return func(any) error { return o.handleFocusSwitch(fc) }
}

func (o *SettingsOverlay) handleFocusSwitch(fc focusedComponent) (err error) {
	if o.focusedComponent == fc {
		return
	}

	if o.focusedComponent != nil {
		err = o.deFocusComponent()
	}
	o.focusedComponent = fc
	o.focusedComponent.Focus()
	return
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
