package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"vbz/orgb"
)

type VBZ struct {
	conn         *orgb.ORGBConn
	countrollers []orgb.Controller
}

func initVBZ() (VBZ, error) {
	conn, err := orgb.Connect("localhost", 6742)
	if err != nil {
		return VBZ{}, err
	}

	count, err := conn.GetControllerCount()
	if err != nil {
		return VBZ{}, err
	}

	controllers := make([]orgb.Controller, count)
	for i := 0; i < count; i++ {
		controller, err := conn.GetController(i)
		if err != nil {
			return VBZ{}, err
		}
		controllers[i] = controller
	}

	return VBZ{
		conn:         conn,
		countrollers: controllers,
	}, nil
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

	if _, err := tea.NewProgram(vbz).Run(); err != nil {
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
	}
	return v, nil
}

func (v VBZ) View() string {
	return "hi"
}
