package circle

import (
	"math"
	"vbz/ui/uiData"

	tea "github.com/charmbracelet/bubbletea"
	lb "github.com/crolbar/lipbalm"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
	lbl "github.com/crolbar/lipbalm/layout"
)

type Circle struct {
	d uiData.UiData
}

func Init(d uiData.UiData) *Circle {
	return &Circle{d: d}
}

func (c *Circle) Resize(tea.WindowSizeMsg) {}
func (c *Circle) Update(msg tea.Msg)       {}

func (c Circle) Render(fb *lbfb.FrameBuffer) {
	c.RenderIn(fb, fb.Size())
}

func (c Circle) RenderIn(fb *lbfb.FrameBuffer, rect lbl.Rect) {
	bins := c.d.Fft.Bins
	sum := 0.0
	n := float64(len(c.d.Fft.Bins)) / 3.0
	for i := 0; i < int(n); i++ {
		sum += bins[i]
	}

	amp := sum / n

	var (
		w = int(rect.Width)
		h = int(rect.Height)

		cy = (h / 2) + int(rect.Y)
		cx = (w / 2) + int(rect.X)

		r  = float64(h/2) - amp
		r2 = float64(h/2-7) - amp
	)

	stepSize := 0.01

	ringChars := []string{"•", "◦", "○", "◎", "●"}

	for i := 0.0; i < 2*math.Pi; i += stepSize {
		x := cx + int((r*amp)*math.Cos(i)*2.2)
		y := cy + int((r*amp)*math.Sin(i))

		charIndex := int(math.Sin(i+float64(c.d.TickCount*3))*2+2) % len(ringChars)
		colorIndex := int(i*30+i*20) + int(c.d.TickCount)*3
		color := lb.Color(uint8(16 + (colorIndex % 216)))

		if x < 0 && x >= w && y < 0 && y >= h {
			continue
		}

		fb.RenderString(lb.SetColor(color, ringChars[charIndex]), lbl.NewRect(uint16(x), uint16(y), 1, 1))
	}

	for i := 0.0; i < 2*math.Pi; i += stepSize {
		x := cx + int((r2*amp)*math.Cos(i)*2.2)
		y := cy + int((r2*amp)*math.Sin(i))

		charIndex := int(math.Sin(i+float64(int(c.d.TickCount)*-2))*4+4) % len(ringChars)
		colorIndex := int(i*35+i*15) + int(c.d.TickCount)*-2
		color := lb.Color(uint8(16 + (colorIndex % 216)))

		if x < 0 && x >= w && y < 0 && y >= h {
			continue
		}

		fb.RenderString(lb.SetColor(color, ringChars[charIndex]), lbl.NewRect(uint16(x), uint16(y), 1, 1))
	}
}
