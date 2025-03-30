package audioCapture

import (
	"fmt"
	"log"

	"github.com/gen2brain/malgo"
)

type AudioCapture struct {
	Dev *malgo.Device
}

func getDevices() (*malgo.AllocatedContext, []malgo.DeviceInfo, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		// fmt.Printf("LOG> %v\n", message)
	})

	if err != nil {
		log.Fatalf("Failed to initialize context: %v", err)
	}

	devices, err := ctx.Devices(malgo.Capture)
	return ctx, devices, err
}

func InitDevice(devIdx int, bufferSize int, sampleRate int, cb malgo.DataProc) (AudioCapture, error) {
	ctx, devices, err := getDevices()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

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
		log.Fatalf("Failed to init device: %v", err)
	}

	return AudioCapture{
		Dev: device,
	}, nil
}

func (a *AudioCapture) StartDev() error {
	return a.Dev.Start()
}

func PrintDevices() error {
	_, devices, err := getDevices()
	if err != nil {
		return err
	}

	for i, d := range devices {
		fmt.Println(i, d.Name())
	}

	return nil
}
