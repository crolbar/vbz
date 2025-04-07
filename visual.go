package main

import (
	"fmt"
	"math"
	"strings"
	"vbz/fft"

	lb "github.com/crolbar/lipbalm"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
	lbl "github.com/crolbar/lipbalm/layout"
)

func (v VBZ) View() string {
	if v.width == 0 || v.height == 0 {
		return ""
	}
	if len(v.fft.Bins) == 0 {
		return "zero len"
	}

	// return v.renderCircle()
	return v.renderBins()
}

func (v VBZ) renderCircle() string {
	bins := v.fft.Bins
	sum := 0.0
	n := float64(len(v.fft.Bins)) / 3.0
	for i := 0; i < int(n); i++ {
		sum += bins[i]
	}

	amp := sum / n

	var (
		w  = v.width
		h  = v.height
		fb = lbfb.NewFrameBuffer(uint16(w), uint16(h))

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

		charIndex := int(math.Sin(i+float64(v.tickCount*3))*2+2) % len(ringChars)
		colorIndex := int(i*30+i*20) + int(v.tickCount)*3
		color := lb.Color(uint8(16 + (colorIndex % 216)))

		if x < 0 && x >= w && y < 0 && y >= h {
			continue
		}

		fb.RenderString(lb.SetColor(color, ringChars[charIndex]), lbl.NewRect(uint16(x), uint16(y), 1, 1))
	}

	for i := 0.0; i < 2*math.Pi; i += stepSize {
		x := cx + int((r2*amp)*math.Cos(i)*2.2)
		y := cy + int((r2*amp)*math.Sin(i))

		charIndex := int(math.Sin(i+float64(int(v.tickCount)*-2))*4+4) % len(ringChars)
		colorIndex := int(i*35+i*15) + int(v.tickCount)*-2
		color := lb.Color(uint8(16 + (colorIndex % 216)))

		if x < 0 && x >= w && y < 0 && y >= h {
			continue
		}

		fb.RenderString(lb.SetColor(color, ringChars[charIndex]), lbl.NewRect(uint16(x), uint16(y), 1, 1))
	}

	return fb.View()
}

func (v VBZ) renderBins() string {
	var (
		bins = v.fft.Bins

		w = v.width
		h = v.height

		fb = lbfb.NewFrameBuffer(uint16(w), uint16(h))

		binWidth     = int(float32(w) / float32(fft.BINS_SIZE))
		binMaxHeight = h
		binE         = strings.Repeat(" ", binWidth)
		binF         = strings.Repeat("█", binWidth)
		startX       = (w - (binWidth * fft.BINS_SIZE)) / 2

		peakLow = v.getPeakLowAmp()
		hueRate = math.Pow(v.hueRate+(v.bpm.bpm*1e-4), 0.99)
	)

	if binWidth == 0 {
		fb.RenderString(
			lb.ExpandHorizontal(w, lb.Center,
				fmt.Sprintf("Min width required: %d, current: %d", fft.BINS_SIZE, w)),
			lbl.NewRect(0, uint16(h/2), uint16(w), 1))
		return fb.View()
	}

	for i, mag := range bins {
		if i >= v.width {
			break
		}
		// TODO ? bpm scaling the hueRate is a bit distracting.
		// maybe add a option to enable/disable the bpm scaling
		v.prevBHues[i] = math.Mod(1+v.prevBHues[i]-v.hueRate, 1)
		v.prevFHues[i] = math.Mod(1+v.prevFHues[i]-hueRate*0.7, 1)

		barHeight := min(int(mag*float64(binMaxHeight)), binMaxHeight)
		bar := strings.Repeat(
			lb.SetColor(
				lb.ColorBgRGB(HSVtoRGB(v.prevBHues[i], 1, peakLow*0.7)),
				binE+"\n",
			), binMaxHeight-barHeight)

		bar = bar + strings.Repeat(
			lb.SetColor(
				lb.ColorRGB(HSVtoRGB(v.prevFHues[i], 1, 1)),
				binF+"\n",
			), barHeight)

		fb.RenderString(bar, lbl.NewRect(uint16(startX+i*binWidth), 0, uint16(binWidth), uint16(h)))
	}

	// fill the black spots at left & right
	if startX != 0 && v.fillBins {
		v.applyFillToBins(&fb, startX, binMaxHeight, binWidth, peakLow, binF, w, h)
	}

	if !v.debug {
		return fb.View()
	}

	// DEBUG

	_v := []string{
		fmt.Sprintf("fps: %d", v.fps),
		fmt.Sprintf("w: %d, h: %d", w, h),
		fmt.Sprintf("bpm: %.2f", v.bpm.bpm),
		fmt.Sprintf("hueRate: %.4f", hueRate),
		fmt.Sprintf("startX: %d", startX),
	}
	for i, str := range _v {
		fb.RenderString(str, lbl.NewRect(0, uint16(i), 15, 1))
	}

	color := lb.ColorBg(0)
	if v.bpm.hasBeat {
		color = lb.ColorBg(1)
	}

	box := "               " + strings.Repeat("\n               ", 6)
	fb.RenderString(lb.SetColor(color, box), lbl.NewRect(uint16(w-15), 0, 15, 5))

	return fb.View()
}

func (v *VBZ) applyFillToBins(
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
		firstHue = v.prevBHues[0]
		lastHue  = v.prevBHues[len(v.prevBHues)-1]

		binHueDiff = (2.0 / float64(fft.BINS_SIZE)) - (1.0 / float64(fft.BINS_SIZE))

		hue = firstHue
	)

	// left
	for i := startX; i >= 0; i-- {
		hue = math.Mod(hue-binHueDiff+1.0, 1)
		fb.RenderString(lb.SetColor(lb.ColorBgRGB(HSVtoRGB(hue, 1, peakLow*0.7)), fill), lbl.NewRect(uint16(i), 0, 1, uint16(h)))

		var (
			barIdx    = startX - i
			mag       = v.fft.Bins[barIdx]
			barHeight = min(int(mag*float64(binMaxHeight)), binMaxHeight)
			bar       = strings.Repeat(lb.SetColor(lb.ColorRGB(HSVtoRGB(v.prevFHues[barIdx], 1, 1)), binF+"\n"), barHeight)
		)

		fb.RenderString(
			lb.SetColor(lb.ColorBgRGB(HSVtoRGB(v.prevFHues[barIdx], 1, 1)), bar),
			lbl.NewRect(uint16(i), uint16(h-barHeight), uint16(binWidth), uint16(barHeight)),
		)
	}

	// right
	hue = lastHue
	for i := w - startX - 1; i < w; i++ {
		hue = math.Mod(hue+binHueDiff+1.0, 1)
		fb.RenderString(lb.SetColor(lb.ColorBgRGB(HSVtoRGB(hue, 1, peakLow*0.7)), fill), lbl.NewRect(uint16(i), 0, 1, uint16(h)))

		var (
			barIdx    = fft.BINS_SIZE - 1 - (i - (w - startX - 1))
			mag       = v.fft.Bins[barIdx]
			barHeight = min(int(mag*float64(binMaxHeight)), binMaxHeight)
			bar       = strings.Repeat(lb.SetColor(lb.ColorRGB(HSVtoRGB(v.prevFHues[barIdx], 1, 1)), binF+"\n"), barHeight)
		)

		fb.RenderString(
			lb.SetColor(lb.ColorBgRGB(HSVtoRGB(v.prevFHues[barIdx], 1, 1)), bar),
			lbl.NewRect(uint16(i), uint16(h-barHeight), uint16(binWidth), uint16(barHeight)),
		)
	}
}
