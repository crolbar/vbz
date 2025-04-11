package ui

import (
	"fmt"
	"math"
	"strings"
	"vbz/fft"
	"vbz/hues"

	lb "github.com/crolbar/lipbalm"
	lbl "github.com/crolbar/lipbalm/layout"
)

func (ui Ui) renderBins() string {
	var (
		bins = ui.fft.Bins

		w = ui.Width
		h = ui.Height

		binWidth     = int(float32(w) / float32(fft.BINS_SIZE))
		binMaxHeight = h
		binE         = strings.Repeat(" ", binWidth)
		binF         = strings.Repeat("█", binWidth)
		startX       = (w - (binWidth * fft.BINS_SIZE)) / 2

		peakLow = ui.fft.PeakLowAmp
	)

	if binWidth == 0 {
		ui.fb.RenderString(
			lb.ExpandHorizontal(w, lb.Center,
				fmt.Sprintf("Min width required: %d, current: %d", fft.BINS_SIZE, w)),
			lbl.NewRect(0, uint16(h/2), uint16(w), 1))
		return ui.fb.View()
	}

	for i, mag := range bins {
		if i >= ui.Width {
			break
		}

		barHeight := min(int(mag*float64(binMaxHeight)), binMaxHeight)
		bar := strings.Repeat(
			lb.SetColor(
				lb.ColorBgRGB(hues.HSVtoRGB(ui.hues.PrevBHues[i], 1, peakLow*0.7)),
				binE+"\n",
			), binMaxHeight-barHeight)

		bar = bar + strings.Repeat(
			lb.SetColor(
				lb.ColorRGB(hues.HSVtoRGB(ui.hues.PrevFHues[i], 1, 1)),
				binF+"\n",
			), barHeight)

		ui.fb.RenderString(bar, lbl.NewRect(uint16(startX+i*binWidth), 0, uint16(binWidth), uint16(h)))
	}

	// fill the black spots at left & right
	if startX != 0 && ui.sets.FillBins {
		ui.applyFillToBins(startX, binMaxHeight, binWidth, peakLow, binF, w, h)
	}

	if !ui.sets.Debug {
		return ui.fb.View()
	}

	// DEBUG

	_v := []string{
		fmt.Sprintf("fps: %d", ui.FPS),
		fmt.Sprintf("w: %d, h: %d", w, h),
		fmt.Sprintf("port: %d", ui.sets.Port),
		fmt.Sprintf("dev: %d", ui.sets.DeviceIdx),
		fmt.Sprintf("filterMode: %d", ui.sets.FilterMode),
		fmt.Sprintf("bpm: %.2f", ui.bpm.Bpm),
		fmt.Sprintf("hueRate: %.4f", math.Pow(ui.sets.HueRate+(ui.bpm.Bpm*1e-4), 0.99)),
		// fmt.Sprintf("click: %v", *v.click),
		// fmt.Sprintf("key: %v", v.key),
		// fmt.Sprintf("mx: %d my: %d", v.mouse.X, v.mouse.Y),
	}
	for i, str := range _v {
		ui.fb.RenderString(str, lbl.NewRect(0, uint16(i), 15, 1))
	}

	color := lb.ColorBg(0)
	if ui.bpm.HasBeat {
		color = lb.ColorBg(1)
	}

	box := "               " + strings.Repeat("\n               ", 6)
	ui.fb.RenderString(lb.SetColor(color, box), lbl.NewRect(uint16(w-15), 0, 15, 5))

	return ui.fb.View()
}

func (ui Ui) applyFillToBins(
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
		firstHue = ui.hues.PrevBHues[0]
		lastHue  = ui.hues.PrevBHues[len(ui.hues.PrevBHues)-1]

		binHueDiff = (2.0 / float64(fft.BINS_SIZE)) - (1.0 / float64(fft.BINS_SIZE))

		hue = firstHue
	)

	// left
	for i := startX; i >= 0; i-- {
		hue = math.Mod(hue-binHueDiff+1.0, 1)
		ui.fb.RenderString(lb.SetColor(lb.ColorBgRGB(hues.HSVtoRGB(hue, 1, peakLow*0.7)), fill), lbl.NewRect(uint16(i), 0, 1, uint16(h)))

		var (
			barIdx    = startX - i
			mag       = ui.fft.Bins[barIdx]
			barHeight = min(int(mag*float64(binMaxHeight)), binMaxHeight)
			bar       = strings.Repeat(lb.SetColor(lb.ColorRGB(hues.HSVtoRGB(ui.hues.PrevFHues[barIdx], 1, 1)), binF+"\n"), barHeight)
		)

		ui.fb.RenderString(
			lb.SetColor(lb.ColorBgRGB(hues.HSVtoRGB(ui.hues.PrevFHues[barIdx], 1, 1)), bar),
			lbl.NewRect(uint16(i), uint16(h-barHeight), uint16(binWidth), uint16(barHeight)),
		)
	}

	// right
	hue = lastHue
	for i := w - startX - 1; i < w; i++ {
		hue = math.Mod(hue+binHueDiff+1.0, 1)
		ui.fb.RenderString(lb.SetColor(lb.ColorBgRGB(hues.HSVtoRGB(hue, 1, peakLow*0.7)), fill), lbl.NewRect(uint16(i), 0, 1, uint16(h)))

		var (
			barIdx    = fft.BINS_SIZE - 1 - (i - (w - startX - 1))
			mag       = ui.fft.Bins[barIdx]
			barHeight = min(int(mag*float64(binMaxHeight)), binMaxHeight)
			bar       = strings.Repeat(lb.SetColor(lb.ColorRGB(hues.HSVtoRGB(ui.hues.PrevFHues[barIdx], 1, 1)), binF+"\n"), barHeight)
		)

		ui.fb.RenderString(
			lb.SetColor(lb.ColorBgRGB(hues.HSVtoRGB(ui.hues.PrevFHues[barIdx], 1, 1)), bar),
			lbl.NewRect(uint16(i), uint16(h-barHeight), uint16(binWidth), uint16(barHeight)),
		)
	}
}

func (ui Ui) renderCircle() string {
	bins := ui.fft.Bins
	sum := 0.0
	n := float64(len(ui.fft.Bins)) / 3.0
	for i := 0; i < int(n); i++ {
		sum += bins[i]
	}

	amp := sum / n

	var (
		w  = ui.Width
		h  = ui.Height

		cy = h / 2
		cx = w / 2

		r  = float64(h/2) - amp
		r2 = float64(h/2-7) - amp
	)

	stepSize := 0.01

	ringChars := []string{"•", "◦", "○", "◎", "●"}

	for i := 0.0; i < 2*math.Pi; i += stepSize {
		x := cx + int((r*amp)*math.Cos(i)*2.2)
		y := cy + int((r*amp)*math.Sin(i))

		charIndex := int(math.Sin(i+float64(ui.TickCount*3))*2+2) % len(ringChars)
		colorIndex := int(i*30+i*20) + int(ui.TickCount)*3
		color := lb.Color(uint8(16 + (colorIndex % 216)))

		if x < 0 && x >= w && y < 0 && y >= h {
			continue
		}

		ui.fb.RenderString(lb.SetColor(color, ringChars[charIndex]), lbl.NewRect(uint16(x), uint16(y), 1, 1))
	}

	for i := 0.0; i < 2*math.Pi; i += stepSize {
		x := cx + int((r2*amp)*math.Cos(i)*2.2)
		y := cy + int((r2*amp)*math.Sin(i))

		charIndex := int(math.Sin(i+float64(int(ui.TickCount)*-2))*4+4) % len(ringChars)
		colorIndex := int(i*35+i*15) + int(ui.TickCount)*-2
		color := lb.Color(uint8(16 + (colorIndex % 216)))

		if x < 0 && x >= w && y < 0 && y >= h {
			continue
		}

		ui.fb.RenderString(lb.SetColor(color, ringChars[charIndex]), lbl.NewRect(uint16(x), uint16(y), 1, 1))
	}

	return ui.fb.View()
}
