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

var l = lbl.DefaultLayout()

func (m Master) Render(fb *lbfb.FrameBuffer) {
	var (
		rows = l.Vercital().
			Constrains(
				lbl.NewConstrain(lbl.Length, 1),
				lbl.NewConstrain(lbl.Length, 1),
				lbl.NewConstrain(lbl.Percent, 65),
				lbl.NewConstrain(lbl.Percent, 35),
			).Split(fb.Size())

		cols = l.Horizontal().
			Constrains(
				lbl.NewConstrain(lbl.Min, 0),
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Min, 0),
				lbl.NewConstrain(lbl.Percent, 50),
				lbl.NewConstrain(lbl.Min, 0),
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Min, 0),
			).Split(rows[2])
	)

	fb.RenderString(
		lb.SetColor(lb.Color(1), lb.ExpandHorizontal(int(fb.Size().Width), lb.Center, "master")),
		rows[0],
	)

	fb.RenderString(m.d.Sets.Host, rows[1], lb.Center)

	m.c.RenderIn(fb, cols[3])
	m.b.RenderIn(fb, rows[3])
}

func (m Master) RenderIn(fb *lbfb.FrameBuffer, rect lbl.Rect) {}
