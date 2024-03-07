package elevator

import (
	"ElevatorLib/elevio"
)

type ElevatorBehaviour string

const (
	EB_Idle     ElevatorBehaviour = "idle"
	EB_DoorOpen ElevatorBehaviour = "door open"
	EB_Moving   ElevatorBehaviour = "moving"
)

type Elevator struct {
	Id        int
	Floor     int
	Dirn      elevio.MotorDirection
	Requests  [4][3]bool
	Behaviour ElevatorBehaviour
}

func Elevator_init(id int) Elevator {
	return Elevator{
		Id:        id,
		Floor:     -1,
		Dirn:      elevio.MD_Stop,
		Behaviour: EB_Idle,
	}
}
