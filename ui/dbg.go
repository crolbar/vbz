package ui

import (
	"fmt"
	"math"
	"strings"

	lb "github.com/crolbar/lipbalm"
	lbl "github.com/crolbar/lipbalm/layout"
)

func (ui Ui) renderDebug() {
	var (
		fbsize = ui.fb.Size()
		w      = int(fbsize.Width)
		h      = int(fbsize.Height)
	)

	_v := []string{
		fmt.Sprintf("fps: %d", ui.d.FPS),
		fmt.Sprintf("w: %d, h: %d", w, h),
		fmt.Sprintf("port: %d", ui.d.Sets.Port),
		fmt.Sprintf("dev: %d", ui.d.Sets.DeviceIdx),
		fmt.Sprintf("filterMode: %d", ui.d.Sets.FilterMode),
		fmt.Sprintf("bpm: %.2f", ui.d.Bpm.Bpm),
		fmt.Sprintf("hueRate: %.4f", math.Pow(ui.d.Sets.HueRate+(ui.d.Bpm.Bpm*1e-4), 0.99)),
		fmt.Sprintf("ampScalar: %d", ui.d.Sets.AmpScalar),
	}
	for i, str := range _v {
		ui.fb.RenderString(str, lbl.NewRect(0, uint16(i), 15, 1))
	}

	color := lb.ColorBg(0)
	if ui.d.Bpm.HasBeat {
		color = lb.ColorBg(1)
	}

	box := "               " + strings.Repeat("\n               ", 6)

	ui.fb.RenderString(lb.SetColor(color, box), lbl.NewRect(uint16(w-15), 0, 15, 5))
}
