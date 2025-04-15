package led

import (
	"vbz/hues"
	"vbz/orgb"
)

type LED struct {
	Conn         *orgb.ORGBConn
	Countrollers []orgb.Controller
}

func InitLED(host string, port int) (*LED, error) {
	var led LED

	conn, err := orgb.Connect(host, port)
	if err != nil {
		return &led, err
	}
	led.Conn = conn

	count, err := conn.GetControllerCount()
	if err != nil {
		return &led, err
	}

	led.Countrollers = make([]orgb.Controller, count)
	for i := 0; i < count; i++ {
		controller, err := conn.GetController(i)
		if err != nil {
			return &led, err
		}
		led.Countrollers[i] = controller
	}

	return &led, nil
}

func (l *LED) SetVibe(
	h *hues.Hues,
	PeakLowAmp float64,
) {
	hI := 0
	for i := 0; i < len(l.Countrollers); i++ {
		c := l.Countrollers[i]

		colors := make([]orgb.RGBColor, len(c.Colors))

		for j := 0; j < len(colors); j++ {
			r, g, b := hues.HSVtoRGB(h.PrevBHues[hI], 1, PeakLowAmp)
			colors[j] = orgb.RGBColor{Red: r, Green: g, Blue: b}
			hI = (hI + 2) % 64
		}

		l.Conn.UpdateLEDS(i, colors)
	}
}

func (l *LED) TurtOffRGB() {
	l.SetAllLEDsToColor(0, 0, 0)
}

func (l *LED) SetAllLEDsToColor(r, g, b uint8) error {
	for i, c := range l.Countrollers {
		colors := make([]orgb.RGBColor, len(c.Colors))

		for i := 0; i < len(colors); i++ {
			colors[i] = orgb.RGBColor{Red: r, Green: g, Blue: b}
		}

		err := l.Conn.UpdateLEDS(i, colors)
		if err != nil {
			return err
		}
	}

	return nil
}
