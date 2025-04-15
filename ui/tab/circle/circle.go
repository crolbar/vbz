package circle

import (
	"math"
	"vbz/hues"
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

const stepSize = 0.01

var ringChars = []string{"•", "◌", "○", "◎", "●"}

func (c Circle) RenderIn(fb *lbfb.FrameBuffer, rect lbl.Rect) {
	bins := c.d.Fft.Bins
	sum := 0.0
	n := float64(len(c.d.Fft.Bins)) / 3.0
	for i := 0; i < int(n); i++ {
		sum += bins[i]
	}
	amp := sum / n

	lamp := c.d.Fft.PeakLowAmp

	var (
		w = int(rect.Width)
		h = int(rect.Height)

		cy = (h / 2) + int(rect.Y)
		cx = (w / 2) + int(rect.X)

		m = min(w, h)

		r  = float64(m / 2)
		r2 = float64(m/2 - 7)

		tickC = c.d.FrameData.TickCount
	)

	c.renderCircle(fb, cx, cy, r*lamp, int(tickC)*3, rect)
	c.renderCircle(fb, cx, cy, r2*amp, int(tickC)*-2, rect)
}

func (c Circle) renderCircle(
	fb *lbfb.FrameBuffer,
	cx int,
	cy int,
	r float64,
	tickC int,
	rect lbl.Rect,
) {
	var (
		w = int(rect.Width)
		h = int(rect.Height)
	)

	for i := 0.0; i < 2*math.Pi; i += stepSize {
		x := cx + int(r*math.Cos(i)*2)
		y := cy + int(r*math.Sin(i))

		charIndex := int(math.Sin(i+float64(tickC)*(c.d.Sets.HueRate*8))*4+4) % len(ringChars)
		hueIdx := int(math.Abs(math.Sin(i+float64(tickC)*c.d.Sets.HueRate))*63) % 64
		color := lb.ColorRGB(
			hues.HSVtoRGB(
				c.d.Hues.PrevFHues[hueIdx],
				1,
				float64(charIndex)/4+0.3,
			),
		)

		if x < 0 && x >= w && y < 0 && y >= h {
			continue
		}

		fb.RenderString(
			lb.SetColor(color, ringChars[charIndex]),
			lbl.NewRect(uint16(x), uint16(y), 1, 1),
		)
	}
}
