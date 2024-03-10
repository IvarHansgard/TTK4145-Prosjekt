package elevator

import (
	"elevatorlib/elevio"
)

type ElevatorBehaviour string

const (
	EB_Idle         ElevatorBehaviour = "idle"
	EB_DoorOpen     ElevatorBehaviour = "doorOpen"
	EB_Moving       ElevatorBehaviour = "moving"
	EB_Disconnected ElevatorBehaviour = "disconnected"
)

type Elevator struct {
	Id       int
	Floor    int
	Dirn     elevio.MotorDirection
	Requests [4][3]bool
	//hallRequests [4][2]bool
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
