package main

import "vbz/orgb"

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
