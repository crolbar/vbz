package main

import (
	"fmt"
	"time"
	"vbz/audioCapture"
	"vbz/fft"
	"vbz/orgb"
	"vbz/settings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/malgo"
)

type VBZ struct {
	conn         *orgb.ORGBConn
	countrollers []orgb.Controller
	audio        *audioCapture.AudioCapture
	settings     settings.Settings

	width  int
	height int

	tickCount    uint
	fps          int
	lastTickTime time.Time

	fft *fft.FFT
}

const SampleRate float64 = 10000

var p *tea.Program

type Refresh struct{}

func (v *VBZ) triggerRefresh() {
	p.Send(Refresh{})
}

func (v *VBZ) initORGBConn() error {
	conn, err := orgb.Connect(v.settings.Host, v.settings.Port)
	if err != nil {
		return err
	}
	v.conn = conn

	count, err := conn.GetControllerCount()
	if err != nil {
		return err
	}

	v.countrollers = make([]orgb.Controller, count)
	for i := 0; i < count; i++ {
		controller, err := conn.GetController(i)
		if err != nil {
			return err
		}
		v.countrollers[i] = controller
	}

	return nil
}

func (v *VBZ) initAudio() error {
	audio, err := audioCapture.InitDevice(v.settings.DeviceIdx, fft.BUFFER_SIZE, int(SampleRate), v.onData())
	if err != nil {
		return err
	}

	v.audio = &audio

	return nil
}

func initVBZ() (*VBZ, error) {
	vbz := VBZ{
		tickCount: 0,
	}

	cfgPath, err := getConfigPathFromArgs()
	if err != nil {
		return &VBZ{}, err
	}

	if cfgPath == "" {
		vbz.settings = settings.GetSettingsDefaultPath()
	} else {
		vbz.settings = settings.GetSettings(cfgPath)
	}

	err = vbz.initORGBConn()
	if err != nil {
		return &VBZ{}, err
	}

	err = vbz.initAudio()
	if err != nil {
		return &VBZ{}, err
	}

	vbz.fft = fft.InitFFT(8000, fft.DoubleBoxFilter, 2, 0.2, 80)

	return &vbz, nil
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

	p = tea.NewProgram(vbz, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		return
	}
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
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
	case Refresh:
	}

	v.tickCount++
	now := time.Now()
	frameTime := now.Sub(v.lastTickTime)
	v.lastTickTime = now
	v.fps = int(time.Second / frameTime)

	return v, nil
}

func (v *VBZ) setBins(samples []uint8) {
	v.fft.ApplyFFT(samples)
	v.triggerRefresh()
}

func (v *VBZ) onData() malgo.DataProc {
	return func(pOutputSample, pInputSamples []byte, frameCount uint32) {
		samples, _ := byteToU8(pInputSamples)
		v.setBins(samples)
		// v.setVibe()

		time.Sleep(1 * time.Millisecond)
	}
}
