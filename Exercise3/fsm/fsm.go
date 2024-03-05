package fsm

import (
	"fmt"
	"time"
"os"
	"ElevatorLib/network/bcast"
	"ElevatorLib/network/localip"
	"ElevatorLib/elevator"
	"ElevatorLib/requests"
	"ElevatorLib/driver-go/elevio"
)

func fsm_init(elevator2, elevator3 <- chan Elevator) {
	e1 = elevator.elevator_uninitialized()
	e2 = elevator.elevator_uninitialized()
	e3 = elevator.elevator_uninitialized()

	e2 <- elevator2
	e3 <- elevator3
}

func setAllLights(es Elevator) {
	for f := 0; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			elevio.SetButtonLamp(f, btn, es.requests[f][btn])
		}
	}

}

func fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(MD_Down)
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
	e1.requests[btn_floor][btn_type] = 1
	fsm_send_data(e1)
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
			elevio.setAllLights(e1)
			e1.behaviour = EB_DoorOpen
			time.Sleep(10*time.Second())
			fsm_onDoorTimeout()
		}
	break
default:
	break
}

}
func fsm_onDoorTimeout() {
	switch e1.behaviour {
	case EB_DoorOpen:
		pair := requests_chooseDirection(elevator)
		e1.dirn = pair.dirn
		e1.behaviour = pair.behaviour
		
		switch e1.behaviour{	
		case EB_DoorOpen:	
			time.Sleep(3 * time.second)
			e1.requests = requests_clearAtCurrentFloor(elevator)
			elevio.setAllLights(e1)
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


func fsm_change_to_master(){}

func fsm_recieveData(rx <- chan Elevator){
	e1.requests = rx.requests

}

func fsm_sendData_master(elevatorNum int, elevatorTx <- chan Elevator){
	if elevatorNum == 2{
	elevatorTx <- e2
}else{
	elevatorTx <- e3
}
}


func fsm_sendData_slave(elevatorTx <- chan Elevator){
	elevatorTx <- e1
}

func fsm_run_algo(){
	optimizedQueue := elevator_algo(e1,e2,e3)
	for i := 0; i < 4 i++{
		e1.requests[i][0] = optimizedQueue["one"][i][0]
		e1.requests[i][1] = optimizedQueue["one"][i][1]		
	
		e2.requests[i][0] = optimizedQueue["two"][i][0]
		e2.requests[i][1] = optimizedQueue["two"][i][1]	

		e3.requests[i][0] = optimizedQueue["three"][i][0]
		e3.requests[i][1] = optimizedQueue["three"][i][1]
	}
}