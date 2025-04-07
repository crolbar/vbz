package main

import (
	"math"
	"vbz/orgb"
)

// all h, s and v are values between 0 and 1
func HSVtoRGB(h, s, v float64) (uint8, uint8, uint8) {
	var (
		C = v * s // chroma / range
		M = v     // max rgb value
		m = M - C // min rgb value

		hueAngle = math.Floor(h * 360)
		Hprime   = hueAngle / 60 // split into regions

		u = Hprime - math.Floor(Hprime) // unbound pos in uprising slope

		t = v * s * u       // bound / in range, pos in uprising slope
		q = v * s * (1 - u) // pos in descending slope
	)

	var r, g, b float64

	switch math.Floor(Hprime) {
	case 0:
		r, g, b = M, t, m
	case 1:
		r, g, b = q, M, m
	case 2:
		r, g, b = m, M, t
	case 3:
		r, g, b = m, q, M
	case 4:
		r, g, b = t, m, M
	case 5:
		r, g, b = M, m, q
	}

	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
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

	peak = 1

	t += 0.003
	hue := math.Mod(t, 1.0)

	// scaledPeak := math.Log(1+9*float64(peak)) / math.Log(10)

	r, g, b := HSVtoRGB(hue, 1, peak)

	_ = r
	_ = g
	_ = b

	// v.setAllLEDsToColor(r, g, b)
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
