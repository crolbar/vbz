package audioCapture

import (
	"errors"
	"fmt"
	"vbz/fft"

	"github.com/gen2brain/malgo"
)

type AudioCapture struct {
	Ctx *malgo.AllocatedContext
	Dev *malgo.Device

	NumDevices      int
	FrameDurationMs float64
	SampleRate      float64
	Callback        malgo.DataProc
}

func InitAudioCapture(devIdx int, bufferSize int, sampleRate float64, cb malgo.DataProc) (AudioCapture, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {})

	if err != nil {
		return AudioCapture{}, errors.New(fmt.Sprintf("Failed to initialize context: %v", err))
	}

	device, numDevices, err := InitDevice(ctx, devIdx, bufferSize, sampleRate, cb)
	if err != nil {
		return AudioCapture{}, err
	}

	return AudioCapture{
		Ctx:             ctx,
		Dev:             device,
		SampleRate:      float64(sampleRate),
		FrameDurationMs: (float64(fft.BUFFER_SIZE) / sampleRate) * 1000,
		NumDevices:      numDevices,
		Callback:        cb,
	}, nil
}

func InitDevice(ctx *malgo.AllocatedContext, devIdx int, bufferSize int, sampleRate float64, cb malgo.DataProc) (*malgo.Device, int, error) {
	devices, err := GetDevices(ctx)
	if err != nil {
		return nil, 0, err
	}

	if devIdx >= len(devices) {
		return nil, 0, errors.New(
			fmt.Sprintf("device idx: %d too large", devIdx),
		)
	}
	selDev := devices[devIdx]

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.DeviceID = selDev.ID.Pointer()
	deviceConfig.Capture.Format = malgo.FormatU8
	deviceConfig.Capture.Channels = 1
	deviceConfig.PeriodSizeInFrames = uint32(bufferSize)
	deviceConfig.SampleRate = uint32(sampleRate)

	deviceConfig.PUserData = nil

	captureCallbacks := malgo.DeviceCallbacks{Data: cb}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return nil, 0, errors.New(
			fmt.Sprintf("Failed to init device: %v", err),
		)
	}

	return device, len(devices), nil
}

func (a *AudioCapture) ReinitDevice(devIdx int) error {
	a.Dev.Stop()
	a.Dev.Uninit()

	device, numDevices, err := InitDevice(
		a.Ctx,
		devIdx,
		fft.BUFFER_SIZE,
		a.SampleRate,
		a.Callback)

	a.Dev = device
	a.NumDevices = numDevices

	return err
}

func (a *AudioCapture) StartDev() error {
	return a.Dev.Start()
}

func GetDevices(ctx *malgo.AllocatedContext) ([]malgo.DeviceInfo, error) {
	devices, err := ctx.Devices(malgo.Capture)
	return devices, err
}

func (a *AudioCapture) PrintDevices() error {
	devices, err := GetDevices(a.Ctx)
	if err != nil {
		return err
	}

	for i, d := range devices {
		fmt.Println(i, d.Name())
	}

	return nil
}
