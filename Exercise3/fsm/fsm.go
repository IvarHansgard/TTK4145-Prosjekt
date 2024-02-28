package fsm

import (
	"fmt"
	"time"

	"../driver-go/elevio"
)

func fsm_init() {
	e1 = elevator.elevator_uninitialized()
}

func setAllLights(es Elevator) {
	for f := 0; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			elevio.SetButtonLamp(btn, f, es.requests[f][btn])
		}
	}

}

func fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	e1.dirn = MD_Down
	e1.behaviour = EB_Moving
}
func fsm_onRequestButtonPress_master(btn_floor int, btn_type buttonType, hallRequests chan <- bool) {
	if btn_type == elevio.BT_HallDown || btn_type == elevio.BT_HallUp {
		hallRequests[btn_floor][btn_type] <- 1
	}else{
		e1.requests[btn_floor][btn_type] = 1
	}
}

func fsm_onRequestButtonPress_slave(btn_floor int, btn_type buttonType) {
}



func onFloorArival(newFloor int) {

	e1.floor = newFloor
	elvio.floorIndicatorLight(elevator.floor)
	switch elevator.behaviour {
	case EB_Moving:
		if requests_shouldStop(e1) {
			SetMotorDirection(elevio.MD_Stop)
			set_door_open_lamp(true)
			e1 = requests_clearAtCurrentFloor(e1)
			time.Sleep(3 * time.second)
			setall= e1.elevator_uninitialized()
	}
}
func fsm_onDoorTimeout() {
	switch elevator.behaviour {
	case EB_DoorOpen:
		pair := requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			time.Sleep(3 * time.second)
			elevator = requests_clearAtCurrentFloor(elevator)
			setAllLights(elevator)
			break
		case EB_Moving:
		case EB_idle:
			doorlight(false)
			motorDirection(elevator.dirn)
			break
		}
		break= elevator.elevator_uninitialized()
		break
	}
}
