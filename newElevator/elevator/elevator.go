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
	Id        int
	StrID     string
	Floor     int
	Dirn      elevio.MotorDirection
	Requests  [4][3]bool
	Behaviour ElevatorBehaviour
}

func Elevator_init(id int, strId string) Elevator {
	return Elevator{
		Id:        id,
		StrID:     strId,
		Floor:     0,
		Dirn:      elevio.MD_Stop,
		Behaviour: EB_Idle,
	}
}
func Elevator_init_disconnected(id int, strId string) Elevator {
	return Elevator{
		Id:        id,
		StrID:     strId,
		Floor:     0,
		Dirn:      elevio.MD_Stop,
		Behaviour: EB_Disconnected,
	}
}
