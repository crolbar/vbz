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
	l.SetAllLEDsToColor(hues.HSVtoRGB(h.PrevFHues[0], 1, PeakLowAmp))
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
