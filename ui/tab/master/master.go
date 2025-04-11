package master

import (
	"vbz/ui/tab"
	"vbz/ui/uiData"

	tea "github.com/charmbracelet/bubbletea"
	lb "github.com/crolbar/lipbalm"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
	lbl "github.com/crolbar/lipbalm/layout"
)

type Master struct {
	d uiData.UiData
	b tab.Tab
	c tab.Tab
}

func Init(d uiData.UiData, b tab.Tab, c tab.Tab) *Master {
	return &Master{
		d: d,
		b: b,
		c: c,
	}
}

func (m *Master) Resize(tea.WindowSizeMsg) {
}

func (m *Master) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	}

}

func (m Master) Render(fb *lbfb.FrameBuffer) {
	m.b.RenderIn(fb, fb.Size())
	m.c.RenderIn(fb, fb.Size())

	fb.RenderString(
		lb.SetColor(lb.Color(1), lb.ExpandHorizontal(int(fb.Size().Width), lb.Center, "master")),
		lbl.NewRect(0, 0, fb.Size().Width, 1),
	)

	var (
	// v = lbl.DefaultLayout().Vercital().
	// 	Constrains(
	// 		lbl.NewConstrain(lbl.Percent, 50),
	// 		lbl.NewConstrain(lbl.Percent, 50),
	// 	).Split(fb.Size())
	)
}

func (m Master) RenderIn(fb *lbfb.FrameBuffer, rect lbl.Rect) {}
