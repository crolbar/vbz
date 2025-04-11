package tab

import (
	tea "github.com/charmbracelet/bubbletea"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
	lbl "github.com/crolbar/lipbalm/layout"
)

type TabType int

const (
	Master TabType = iota
	Bins
	Circle
	Last__
)

type Tab interface {
	Resize(tea.WindowSizeMsg)
	Update(tea.Msg)
	Render(*lbfb.FrameBuffer)
	RenderIn(*lbfb.FrameBuffer, lbl.Rect)
}
