package ui

import "vbz/settings"

func assertSettings(t any) *settings.Settings {
	var (
		s  *settings.Settings
		ok bool
	)
	if s, ok = t.(*settings.Settings); !ok {
		panic(ok)
	}

	return s
}
