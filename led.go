package main

import (
	"math"
	"vbz/orgb"
)

func HSVtoRGB(h, s, v float64) (float64, float64, float64) {
	i := int(h * 6)
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	i = i % 6

	switch i {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	case 5:
		return v, p, q
	}
	return 0, 0, 0
}

var t float64

func (v *VBZ) setVibe() {
	peak := 0.0
	// n := len(*v.bins)/3
	n := 1
	bins := (v.fft.Bins)
	for i := 0; i < int(n); i++ {
		if peak < bins[i] {
			peak = bins[i]
		}
	}

	t += 0.01

	hue := math.Mod(t, 1.0)

	// smooth out peak
	scaledPeak := math.Log(1+9*float64(peak)) / math.Log(10)
	r, g, b := HSVtoRGB(hue, 1.0, scaledPeak)

	v.setAllLEDsToColor(uint8(r*255), uint8(g*255), uint8(b*255))
}

func (v *VBZ) turtOffRGB() {
	v.setAllLEDsToColor(0, 0, 0)
}

func (v *VBZ) setAllLEDsToColor(r, g, b uint8) error {
	for i, c := range v.countrollers {
		colors := make([]orgb.RGBColor, len(c.Colors))

		for i := 0; i < len(colors); i++ {
			colors[i] = orgb.RGBColor{Red: r, Green: g, Blue: b}
		}

		err := v.conn.UpdateLEDS(i, colors)
		if err != nil {
			return err
		}
	}

	return nil
}
