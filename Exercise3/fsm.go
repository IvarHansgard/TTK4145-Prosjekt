package main

import (
	"fmt"
	"time"
)

func fsm_init() {
	elevator = elevator_uninitialized()
	outputDeive = getOutputDevice()
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
	elevator.dirn = MD_Down
	elevator.behaviour = EB_Moving
}
func fsm_onRequestButtonPress(btn_floor int, btn_type buttonType) {
	fmt.Println(btn_floor, eb_toString(btn_type))
	elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requests_shouldClearImmediately(elevator, btn_floor, btn_type) {
			time.sleep(3 * time.second)
		} else {
			elevator.requests[btn_floor][btn_type] = true
		}
		break

	case EB_Moving:
		elevator.requests[btn_floor][btn_type] = true
		break

	case EB_idle:
		elevator.requests[btn_floor][btn_type] = true
		pair := requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch pair.behaviour {
		case EB_DoorOpen:
			doorlight(true)
			time.sleep(3 * time.second)
			elevator = requests_clearAtCurrentFloor(elevator)
			break

		case EB_Moving:
			outputDevice.setMotorDirection(elevator.dirn)
			break

		case EB_idle:
			break
		}
		break
	}

	setAllLights(elevator)
	elevator_print(elevator)

}
func fsm_onFloorArival(newFloor int) {
	elevator_print(elevator)
	elevator.floor = newFloor
	elvio.floorIndicatorLight(elevator.floor)
	switch elevator.behaviour {
	case EB_Moving:
		if requests_shouldStop(elevator) {
			SetMotorDirection(elevio.MD_Stop)
			set_door_open_lamp(true)
			elevator = requests_clearAtCurrentFloor(elevator)
			time.Sleep(3 * time.second)
			setallLights(elevator)
			elevator.behaviour = EB_DoorOpen
		}
		break
	default:
		break
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
		break
	default:
		break
	}
}
