package settingsOverlay

import lbti "github.com/crolbar/lipbalm/components/textInput"

type focusedComponent interface {
	Focus()
	DeFocus()
	HasFocus() bool
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

func (o *SettingsOverlay) deFocusComponent() {
	if c, ok := o.focusedComponent.(*lbti.TextInput); ok {
		o.handleTiDefocus(c)
	}

	if o.focusedComponent != nil {
		o.focusedComponent.DeFocus()
		o.focusedComponent = nil
	}
}
