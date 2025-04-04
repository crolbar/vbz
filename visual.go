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
	bins := v.fft.Bins

	var (
		w = v.width
		h = v.height

		fb = lbfb.NewFrameBuffer(uint16(w), uint16(h))

		binWidth     = int(float32(w)/float32(fft.BINS_SIZE)) + 2
		binMaxHeight = h
		// binWidth = 3
		// binMaxHeight = 30
		startX = (w - (binWidth * fft.BINS_SIZE)) / 2
	)

	for i, mag := range bins {
		if i >= v.width {
			break
		}

		barHeight := min(int(mag*float64(binMaxHeight)), binMaxHeight)
		bar := strings.Repeat(" \n", binMaxHeight-barHeight)
		bar = bar + strings.Repeat("█\n", barHeight)

		for range binWidth {
			bar = lb.JoinHorizontal(lb.Left, bar, bar)
		}

		fb.RenderString(bar, lbl.NewRect(uint16(startX+i*binWidth), 0, uint16(binWidth), uint16(v.height)))
	}

	fb.RenderString(fmt.Sprintf("fps: %d", v.fps), lbl.NewRect(0, 0, 15, 1))

	return fb.View()
}
