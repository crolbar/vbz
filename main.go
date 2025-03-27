package main

import (
	"fmt"
	"math"
	"time"
	"vbz/audio"
	"vbz/orgb"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/malgo"
)

type VBZ struct {
	conn         *orgb.ORGBConn
	countrollers []orgb.Controller
	audio        *audio.Audio
}

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
func (v *VBZ) setVibe(peak float32) {
    t += 0.01

    hue := math.Mod(t, 1.0)

	// smooth out peak
	scaledPeak := math.Log(1 + 9*float64(peak)) / math.Log(10)
    r, g, b := HSVtoRGB(hue, 1.0, scaledPeak)

    v.setAllLEDsToColor(uint8(r*255), uint8(g*255), uint8(b*255))
}

func (v *VBZ) onData() malgo.DataProc {
	return func(pOutputSample, pInputSamples []byte, frameCount uint32) {
		samples, _ := byteToFloat32(pInputSamples)

		var peak float32 = 0.0
		for _, s := range samples {
			if peak < s {
				peak = s
			}
		}

		fmt.Println(peak)

		v.setVibe(peak)

		time.Sleep(5 * time.Millisecond)
	}
}

func initVBZ() (VBZ, error) {
	vbz := VBZ{}

	conn, err := orgb.Connect("localhost", 6742)
	if err != nil {
		return VBZ{}, err
	}
	vbz.conn = conn

	count, err := conn.GetControllerCount()
	if err != nil {
		return VBZ{}, err
	}

	vbz.countrollers = make([]orgb.Controller, count)
	for i := 0; i < count; i++ {
		controller, err := conn.GetController(i)
		if err != nil {
			return VBZ{}, err
		}
		vbz.countrollers[i] = controller
	}

	audio, err := audio.InitDevice(1, vbz.onData())
	if err != nil {
		return VBZ{}, err
	}

	vbz.audio = &audio

	return vbz, nil
}

func main() {
	vbz, err := initVBZ()
	if err != nil {
		fmt.Println("Error while connecting to openrgb: ", err)
		return
	}

	defer vbz.conn.Close()

	b, err := vbz.parseArgs()
	if err != nil {
		fmt.Println("Error while parsing args: ", err)
		return
	}
	if b {
		return
	}

	vbz.audio.StartDev()
	defer vbz.audio.Dev.Uninit()

	// if _, err := tea.NewProgram(vbz).Run(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	select {}
}

func (v VBZ) Init() tea.Cmd {
	return nil
}

func (v VBZ) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			v.setAllLEDsToColor(255, 0, 0)
		case "g":
			v.setAllLEDsToColor(0, 255, 0)
		case "b":
			v.setAllLEDsToColor(0, 0, 255)
		case "q":
			return v, tea.Quit
		}
	}
	return v, nil
}

func (v VBZ) View() string {
	return "hi"
}
