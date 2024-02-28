package main


import{
	"fmt"
	"elevio"
}

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

// Dirn represents the direction of the elevator.
type Dirn int

const (
	D_Up Dirn = iota
	D_Down
	D_Stop
)

// vet ikke om jeg skal slå sammen med elvator_io_types:
/*
const (
	D_Up        = 1
	D_Down Dirn = -1
	D_Stop      = 0
	
)

*/

// Button represents the different types of buttons in the elevator.
type Button int

const (
	B_HallUp Button = iota
	B_HallDown
	B_Cab
)

func elevio_dirn_toString(d Dirn) string{
	switch d {
	case D_Up:
		return "D_Up"
	case D_Stop:
		return "D_Stop"
	case D_Down:
		return "D_Down"
	}
}

func elevio_button_toString(b Button) string{
	switch d {
	case B_HallUp:
		return "B_HallUp"
	case B_Cab:
		return "B_Cab"
	case B_HallDown:
		return "B_HallDown"
	}
}
type ElevInputDevice struct {
	floorSensor func() int
	requestButton func(Button,int) int
	stopButton func() int
	obstruction func() int
}

ElevOutputDevice struct{
	floorIndicator func() int
	requestButton func(Button,int) int
	doorLight func() int 
	stopButtonLight func() int
	motorDirection func() int
}



//vet ikke om jeg skal implementere og initialisere dette her eller i main, vet heller ikke om det er riktig

inputDevice := ElevInputDevice{
	FloorSensor:   elevio.GetFloor,
	RequestButton: elevio.GetButton,
	StopButton:    elevio.GetStop,
	Obstruction:   elevio.GetObstruction,
}

outputDevice := ElevOutputDevice{
	FloorIndicator:     elevio.SetFloorIndicator,
	RequestButtonLight: elevio.SetButtonLamp,
	DoorLight:          elevio.SetDoorOpenLamp,
	StopButtonLight:    elevio.SetStopLamp,
	MotorDirection:     elevio.SetMotorDirection,
}
