package requests

import (
	"elevatorlib/elevator"
	. "elevatorlib/elevator"
	"elevatorlib/elevio"
	"fmt"
)

type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection
	Behaviour ElevatorBehaviour
}

func CabRequestsAbove(e Elevator) bool {
	for f := e.Floor + 1; f < 4; f++ {
		for btn := 0; btn < 3; btn++ {
			if e.Requests[f][btn] {
				fmt.Println("CabRequestsAbove")
				return true
			}
		}
	}
	return false
}

func CabRequestsBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < 3; btn++ {
			if e.Requests[f][btn] {
				fmt.Println("CabRequestsBelow")
				return true
			}
		}
	}
	return false
}

func CabRequestsHere(e Elevator) bool {
	for btn := 0; btn < 3; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func HallUpRequestHere(e Elevator) bool {
	if e.Requests[e.Floor][elevio.BT_HallUp] {
		return true
	} else {
		return false
	}

}

func HallDownRequestHere(e Elevator) bool {
	if e.Requests[e.Floor][elevio.BT_HallDown] {
		return true
	} else {
		return false
	}

}

func CabRequestHere(e Elevator) bool {
	if e.Requests[e.Floor][elevio.BT_Cab] {
		return true
	} else {
		return false
	}

}

func RequestsChooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevio.MD_Up:
		if CabRequestsAbove(e) {
			fmt.Println("md up, request above")
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		} else if CabRequestsBelow(e) {
			fmt.Println("md up, request below")
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		} else if CabRequestsHere(e) {
			fmt.Println("md up, request here")
			return DirnBehaviourPair{elevio.MD_Down, EB_DoorOpen}

		} else {
			fmt.Println("md up, idle")
			return DirnBehaviourPair{elevio.MD_Stop, elevator.EB_Idle}
		}
	case elevio.MD_Down:
		if CabRequestsBelow(e) {
			fmt.Println("md down, requests below")
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		} else if CabRequestsHere(e) {
			fmt.Println("md down, requests here")
			return DirnBehaviourPair{elevio.MD_Down, EB_DoorOpen}
		} else if CabRequestsAbove(e) {
			fmt.Println("md down, requests above")
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		} else {
			//fmt.Println("md down, idle")
			return DirnBehaviourPair{elevio.MD_Stop, EB_Idle}
		}

	case elevio.MD_Stop:
		if CabRequestsHere(e) {
			fmt.Println("md stop, requests here")
			return DirnBehaviourPair{elevio.MD_Stop, EB_DoorOpen}
		} else if CabRequestsAbove(e) {
			fmt.Println("md stop, requests above")
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		} else if CabRequestsBelow(e) {
			fmt.Println("md stop, requests below")
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		} else {
			fmt.Println("md stop, idle")
			return DirnBehaviourPair{elevio.MD_Stop, EB_Idle}
		}

	default:
		fmt.Println("default")
		return DirnBehaviourPair{elevio.MD_Stop, EB_Idle}
	}
}

func RequestsShouldStop(e Elevator) bool {
	switch e.Dirn {
	case elevio.MD_Down:
		return HallDownRequestHere(e) || CabRequestHere(e) || !CabRequestsBelow(e) //bool eller int?
	case elevio.MD_Up:
		return HallUpRequestHere(e) || CabRequestHere(e) || !CabRequestsAbove(e)
	case elevio.MD_Stop:
		fallthrough
	default:
		return true
	}
}

func RequestsShouldClearImmediately(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
	return e.Floor == btn_floor && ((e.Dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) || (e.Dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) || (e.Dirn == elevio.MD_Stop || btn_type == elevio.BT_Cab))
}

func RequestsClearAtCurrentFloor(e Elevator) Elevator {

	e.Requests[e.Floor][elevio.BT_Cab] = false
	switch e.Dirn {
	case elevio.MD_Up:
		if !CabRequestsAbove(e) && !HallUpRequestHere(e) {
			e.Requests[e.Floor][elevio.BT_HallDown] = false
		}

		e.Requests[e.Floor][elevio.BT_HallUp] = false

	case elevio.MD_Down:
		if !CabRequestsBelow(e) && !HallDownRequestHere(e) {
			e.Requests[e.Floor][elevio.BT_HallUp] = false
		}
		e.Requests[e.Floor][elevio.BT_HallDown] = false

	case elevio.MD_Stop:
		fallthrough
	default:
		e.Requests[e.Floor][elevio.BT_HallUp] = false
		e.Requests[e.Floor][elevio.BT_HallDown] = false
	}

	return e

}

func HallRequestsClearAtCurrentFloor(e Elevator) elevio.ButtonEvent {
	var buttonToclear elevio.ButtonEvent
	switch e.Dirn {
	case elevio.MD_Down:
		buttonToclear.Floor = e.Floor
		buttonToclear.Button = elevio.BT_HallDown
	case elevio.MD_Up:
		buttonToclear.Floor = e.Floor
		buttonToclear.Button = elevio.BT_HallUp
	case elevio.MD_Stop:
		buttonToclear.Floor = e.Floor
		if HallUpRequestHere(e) && !HallDownRequestHere(e) {
			buttonToclear.Button = elevio.BT_HallUp
		} else if HallDownRequestHere(e) && !HallUpRequestHere(e) {
			buttonToclear.Button = elevio.BT_HallDown
		}

		if HallUpRequestHere(e) && HallDownRequestHere(e) {
			if !CabRequestsBelow(e) {
				buttonToclear.Button = elevio.BT_HallDown
				e.Behaviour = elevator.EB_DoorOpen

			} else if !CabRequestsAbove(e) {
				buttonToclear.Button = elevio.BT_HallUp
				e.Behaviour = elevator.EB_DoorOpen

			}

		}

	}

	if e.Floor == 0 {
		buttonToclear.Floor = e.Floor
		buttonToclear.Button = elevio.BT_HallUp
	} else if e.Floor == 3 {
		buttonToclear.Floor = e.Floor
		buttonToclear.Button = elevio.BT_HallDown
	}
	return buttonToclear

}
