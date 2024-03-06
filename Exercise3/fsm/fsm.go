package fsm

import (
	"ElevatorLib/driver-go/elevio"
	. "ElevatorLib/elevator"
	"time"
)

func Fsm_init() {

}

func SetAllLights(es Elevator) {
	for f := 0; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			elevio.SetButtonLamp(f, btn, es.requests[f][btn])
		}
	}

}

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(MD_Down)
	e1.dirn = MD_Down
	e1.behaviour = EB_Moving
}

func Fsm_onRequestButtonPress_master(btn_floor int, btn_type buttonType, hallRequests chan<- [4][2]int) { //byttet fra [4][2] chan <- int til chan <- [4][2]int
	if btn_type == elevio.BT_HallDown || btn_type == elevio.BT_HallUp {
		hallRequests[btn_floor][btn_type] <- 1
	} else {
		elevatorArray[0].requests[btn_floor][btn_type] = 1
	}
}

func Fsm_onRequestButtonPress_slave(btn_floor int, btn_type buttonType) {
	elevatorArray[0].requests[btn_floor][btn_type] = 1
	fsm_send_data(elevatorArray[0])
}

func OnFloorArrival(newFloor int) {
	elevatorArray[0].floor = newFloor
	elvio.floorIndicatorLight(elevator.floor)
	switch elevator.behaviour {
	case EB_Moving:
		if requests_shouldStop(elevatorArray[0]) {
			SetMotorDirection(elevio.MD_Stop)
			set_door_open_lamp(true)
			elevatorArray[0] = requests_clearAtCurrentFloor(elevatorArray[0])
			time.Sleep(3 * time.second)
			elevio.setAllLights(elevatorArray[0])
			elevatorArray[0].behaviour = EB_DoorOpen
			time.Sleep(10 * time.Second())
			fsm_onDoorTimeout()
		}
		break
	default:
		break
	}

}
func Fsm_onDoorTimeout() {
	switch elevatorArray[0].behaviour {
	case EB_DoorOpen:
		pair := requests_chooseDirection(elevator)
		elevatorArray[0].dirn = pair.dirn
		elevatorArray[0].behaviour = pair.behaviour

		switch elevatorArray[0].behaviour {
		case EB_DoorOpen:
			time.Sleep(3 * time.second)
			elevatorArray[0].requests = requests_clearAtCurrentFloor(elevator)
			elevio.setAllLights(elevatorArray[0])
			break
		case EB_Moving:
		case EB_idle:
			elevio.set_door_open_lamp(false)
			elevio.SetMotorDirection(elevator.dirn)
			break
		}
	}
	break
}
func Fsm_check_master(id int, activeElevators [3]bool) {
	if activeElevators[id-1] == false && id != 0 {
		Fsm_change_to_master()
	}
}
func Fsm_change_to_master() {
	
}

func Fsm_recieveData(rx <-chan Elevator) {
	elevatorArray[0].requests = rx.requests

}

func Fsm_sendData_master(elevatorNum int, elevatorTx <-chan Elevator) {
	if elevatorNum == 2 {
		elevatorTx <- elevatorArray[1]
	} else {
		elevatorTx <- elevatorArray[2]
	}
}

func Fsm_sendData_slave(elevatorTx <-chan Elevator) {
	elevatorTx <- elevatorArray[0]
}

func Fsm_run_algo() {
	optimizedQueue := request_asigner.Elevator_algo(hallRequests, elevatorArray[0], elevatorArray[1], elevatorArray[2])
	for i := 0; i < 4; i++ {
		elevatorArray[0].requests[i][0] = optimizedQueue["one"][i][0]
		elevatorArray[0].requests[i][1] = optimizedQueue["one"][i][1]

		elevatorArray[1].requests[i][0] = optimizedQueue["two"][i][0]
		elevatorArray[1].requests[i][1] = optimizedQueue["two"][i][1]

		elevatorArray[2].requests[i][0] = optimizedQueue["three"][i][0]
		elevatorArray[2].requests[i][1] = optimizedQueue["three"][i][1]
	}
}
