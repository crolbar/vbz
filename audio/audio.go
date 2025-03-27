package audio

import (
	"log"
	"github.com/gen2brain/malgo"
)

type Audio struct {
	Dev *malgo.Device
}

func InitDevice(devIdx int, cb malgo.DataProc) (Audio, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		// fmt.Printf("LOG> %v\n", message)
	})

	if err != nil {
		log.Fatalf("Failed to initialize context: %v", err)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	devices, err := ctx.Devices(malgo.Capture)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("========== Devices ==============")
	// for _, d := range devices {
	// 	fmt.Println(d.Name())
	// }
	// fmt.Println("==========================")

	selDev := devices[devIdx]

	// fmt.Println("selected id: ", selDev.String())

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.DeviceID = selDev.ID.Pointer()
	deviceConfig.Capture.Format = malgo.FormatF32
	deviceConfig.Capture.Channels = 2
	deviceConfig.PUserData = nil

	captureCallbacks := malgo.DeviceCallbacks{Data: cb}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		log.Fatalf("Failed to init device: %v", err)
	}

	return Audio{
		Dev: device,
	}, nil
}

func (a *Audio) StartDev() error {
	return a.Dev.Start()
}
