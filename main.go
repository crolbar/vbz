package main

import (
	"errors"
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

	audio *audioCapture.AudioCapture

	settings   settings.Settings
	configPath string

	width  int
	height int

	shouldNotEnterTui bool

	tickCount    uint
	fps          int
	lastTickTime time.Time

	fft *fft.FFT
}

const SampleRate float64 = 10000

var p *tea.Program

type Refresh struct{}

func (v *VBZ) triggerRefresh() {
	if p != nil {
		p.Send(Refresh{})
	}
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

func initVBZ() (VBZ, error) {
	var err error

	var defaultFFT = fft.DefaultFFT
	vbz := VBZ{
		tickCount: 0,
		fft:       &defaultFFT,
		settings: settings.Settings{
			DeviceIdx: -1,
			Port:      -1,
			Host:      "-1",
			FftPtr:    &defaultFFT,
		},
	}

	err = vbz.parseEarlyArgs()
	if err != nil {
		return VBZ{}, err
	}
	if vbz.shouldNotEnterTui {
		return vbz, nil
	}

	err = vbz.settings.InitSettings(vbz.configPath)
	if err != nil {
		return VBZ{}, errors.New(fmt.Sprintf("Error seting settings: %s", err.Error()))
	}

	err = vbz.initORGBConn()
	if err != nil {
		return VBZ{}, errors.New(fmt.Sprintf("Error initializing openrgb connection: %s", err.Error()))
	}

	err = vbz.initAudio()
	if err != nil {
		return VBZ{}, errors.New(fmt.Sprintf("Error initializing audio capture dev: %s", err.Error()))
	}

	err = vbz.parseLateArgs()
	if err != nil {
		return VBZ{}, errors.New(fmt.Sprintf("Error parsing args: %s", err.Error()))
	}

	return vbz, nil
}

func main() {
	vbz, err := initVBZ()
	if err != nil {
		fmt.Println(err)
		return
	}
	if vbz.shouldNotEnterTui {
		return
	}

	defer vbz.conn.Close()

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
