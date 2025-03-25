package main

import (
	"fmt"
	"log"
	"vbz/orgb"
)

func main() {
	c, err := orgb.Connect("localhost", 6742)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	count, err := c.GetControllerCount()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(count)

	for i := 0; i < count; i++ {
		controller, _ := c.GetController(i)
		fmt.Println(controller)

		colors := make([]orgb.RGBColor, len(controller.Colors))

		for i := 0; i < len(colors); i++ {
			// colors[i] = orgb.Color{uint8(rand.Uint32() % 255), uint8(rand.Uint32() % 255), uint8(rand.Uint32() % 255)}
			colors[i] = orgb.RGBColor{
				// Red:   10,
				// Green: 20,
				// Blue:  10,
				Red:   0,
				Green: 0,
				Blue:  0,
			}
		}

		fmt.Printf("%s\n", controller.Name)
		if err := c.UpdateLEDS(i, colors); err != nil {
			log.Fatal(err)
		}
	}
}
