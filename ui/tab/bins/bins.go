package bins

import (
	"fmt"
	"math"
	"strings"
	"vbz/fft"
	"vbz/hues"
	"vbz/ui/uiData"

	tea "github.com/charmbracelet/bubbletea"
	lb "github.com/crolbar/lipbalm"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
	lbl "github.com/crolbar/lipbalm/layout"
)

type Bins struct {
	d uiData.UiData

	mouse tea.MouseEvent
	key   string
}

func Init(d uiData.UiData) *Bins {
	return &Bins{d: d}
}

func (b *Bins) Resize(tea.WindowSizeMsg) {}
func (b *Bins) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		b.key = msg.String()
	case tea.MouseMsg:
		b.mouse = tea.MouseEvent(msg)
	}
}

func (b Bins) Render(fb *lbfb.FrameBuffer) {
	b.RenderIn(fb, fb.Size())

	fb.RenderString(
		lb.SetColor(lb.Color(2), lb.ExpandHorizontal(int(fb.Size().Width), lb.Center, "bins")),
		lbl.NewRect(0, 0, fb.Size().Width, 1),
	)
}

func (b Bins) RenderIn(fb *lbfb.FrameBuffer, rect lbl.Rect) {
	var (
		bins = b.d.Fft.Bins

		w = int(rect.Width)
		h = int(rect.Height)

		binWidth     = int(float32(w) / float32(fft.BINS_SIZE))
		binMaxHeight = h
		binE         = strings.Repeat(" ", binWidth)
		binF         = strings.Repeat("█", binWidth)
		startX       = (w - (binWidth * fft.BINS_SIZE)) / 2

		peakLow = b.d.Fft.PeakLowAmp
	)

	if binWidth == 0 {
		fb.RenderString(
			lb.ExpandHorizontal(w, lb.Center,
				fmt.Sprintf("Min width required: %d, current: %d", fft.BINS_SIZE, w)),
			lbl.NewRect(0, uint16(h/2), uint16(w), 1))
		return
	}

	for i, mag := range bins {
		if i >= w {
			break
		}

		barHeight := min(int(mag*float64(binMaxHeight)), binMaxHeight)
		bar := strings.Repeat(
			lb.SetColor(
				lb.ColorBgRGB(hues.HSVtoRGB(b.d.Hues.PrevBHues[i], 1, peakLow*0.7)),
				binE+"\n",
			), binMaxHeight-barHeight)

		bar = bar + strings.Repeat(
			lb.SetColor(
				lb.ColorRGB(hues.HSVtoRGB(b.d.Hues.PrevFHues[i], 1, 1)),
				binF+"\n",
			), barHeight)

		fb.RenderString(bar, lbl.NewRect(rect.X+uint16(startX+i*binWidth), rect.Y, uint16(binWidth), uint16(h)))
	}

	// fill the black spots at left & right
	if startX != 0 && b.d.Sets.FillBins {
		b.applyFillToBins(fb, startX, binMaxHeight, binWidth, peakLow, binF, w, h)
	}

	if b.d.Sets.Debug {
		b.renderDebug(fb)
	}
}

func (b Bins) applyFillToBins(
	fb *lbfb.FrameBuffer,
	startX int,
	binMaxHeight int,
	binWidth int,
	peakLow float64,
	binF string,
	w int,
	h int,
) {
	var (
		fill     = strings.Repeat(" \n", h)
		firstHue = b.d.Hues.PrevBHues[0]
		lastHue  = b.d.Hues.PrevBHues[len(b.d.Hues.PrevBHues)-1]

		binHueDiff = (2.0 / float64(fft.BINS_SIZE)) - (1.0 / float64(fft.BINS_SIZE))

		hue = firstHue
	)

	// left
	for i := startX; i >= 0; i-- {
		hue = math.Mod(hue-binHueDiff+1.0, 1)
		fb.RenderString(lb.SetColor(lb.ColorBgRGB(hues.HSVtoRGB(hue, 1, peakLow*0.7)), fill), lbl.NewRect(uint16(i), 0, 1, uint16(h)))

		var (
			barIdx    = startX - i
			mag       = b.d.Fft.Bins[barIdx]
			barHeight = min(int(mag*float64(binMaxHeight)), binMaxHeight)
			bar       = strings.Repeat(lb.SetColor(lb.ColorRGB(hues.HSVtoRGB(b.d.Hues.PrevFHues[barIdx], 1, 1)), binF+"\n"), barHeight)
		)

		fb.RenderString(
			lb.SetColor(lb.ColorBgRGB(hues.HSVtoRGB(b.d.Hues.PrevFHues[barIdx], 1, 1)), bar),
			lbl.NewRect(uint16(i), uint16(h-barHeight), uint16(binWidth), uint16(barHeight)),
		)
	}

	// right
	hue = lastHue
	for i := w - startX - 1; i < w; i++ {
		hue = math.Mod(hue+binHueDiff+1.0, 1)
		fb.RenderString(lb.SetColor(lb.ColorBgRGB(hues.HSVtoRGB(hue, 1, peakLow*0.7)), fill), lbl.NewRect(uint16(i), 0, 1, uint16(h)))

		var (
			barIdx    = fft.BINS_SIZE - 1 - (i - (w - startX - 1))
			mag       = b.d.Fft.Bins[barIdx]
			barHeight = min(int(mag*float64(binMaxHeight)), binMaxHeight)
			bar       = strings.Repeat(lb.SetColor(lb.ColorRGB(hues.HSVtoRGB(b.d.Hues.PrevFHues[barIdx], 1, 1)), binF+"\n"), barHeight)
		)

		fb.RenderString(
			lb.SetColor(lb.ColorBgRGB(hues.HSVtoRGB(b.d.Hues.PrevFHues[barIdx], 1, 1)), bar),
			lbl.NewRect(uint16(i), uint16(h-barHeight), uint16(binWidth), uint16(barHeight)),
		)
	}
}

func (b Bins) renderDebug(fb *lbfb.FrameBuffer) {
	var (
		fbsize = fb.Size()
		w      = int(fbsize.Width)
		h      = int(fbsize.Height)
	)

	_v := []string{
		fmt.Sprintf("fps: %d", b.d.FPS),
		fmt.Sprintf("w: %d, h: %d", w, h),
		fmt.Sprintf("port: %d", b.d.Sets.Port),
		fmt.Sprintf("dev: %d", b.d.Sets.DeviceIdx),
		fmt.Sprintf("filterMode: %d", b.d.Sets.FilterMode),
		fmt.Sprintf("bpm: %.2f", b.d.Bpm.Bpm),
		fmt.Sprintf("hueRate: %.4f", math.Pow(b.d.Sets.HueRate+(b.d.Bpm.Bpm*1e-4), 0.99)),
		fmt.Sprintf("key: %v", b.key),
		fmt.Sprintf("mx: %d my: %d", b.mouse.X, b.mouse.Y),
	}
	for i, str := range _v {
		fb.RenderString(str, lbl.NewRect(0, uint16(i), 15, 1))
	}

	color := lb.ColorBg(0)
	if b.d.Bpm.HasBeat {
		color = lb.ColorBg(1)
	}

	box := "               " + strings.Repeat("\n               ", 6)

	fb.RenderString(lb.SetColor(color, box), lbl.NewRect(uint16(w-15), 0, 15, 5))
}
