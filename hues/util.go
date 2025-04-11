package hues

import "math"

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
