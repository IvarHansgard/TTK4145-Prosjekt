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

// i C
/*
static int CabRequestsAbove(Elevator e){
    for(int f = e.floor+1; f < N_FLOORS; f++){
        for(int btn = 0; btn < N_BUTTONS; btn++){
            if(e.requests[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}*/

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

// i C
/*
static int CabRequestsBelow(Elevator e){
    for(int f = 0; f < e.floor; f++){
        for(int btn = 0; btn < N_BUTTONS; btn++){
            if(e.requests[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}*/

func CabRequestsHere(e Elevator) bool {
	for btn := 0; btn < 3; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func HallUpRequestHere(e Elevator) bool {
	if e.Requests[e.Floor][elevio.BT_HallUp]{
		return true
	}else{
		return false
	}
	

}

func HallDownRequestHere(e Elevator) bool {
	if e.Requests[e.Floor][elevio.BT_HallDown]{
		return true
	}else{
		return false
	}
	

}

func RequestCabHere(e Elevator) bool {
	if  e.Requests[e.Floor][elevio.BT_Cab]{
		return true
	}else{
		return false
	}
	

}

/*
// i C:
static int CabRequestsHere(Elevator e){
    for(int btn = 0; btn < N_BUTTONS; btn++){
        if(e.requests[e.floor][btn]){
            return 1;
        }
    }
    return 0;
}*/

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
			fmt.Println("md down, idle")
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

//i C
/*
DirnBehaviourPair RequestsChooseDirection(Elevator e){
    switch(e.dirn){
    case D_Up:
        return  CabRequestsAbove(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
                CabRequestsHere(e)  ? (DirnBehaviourPair){D_Down, EB_DoorOpen} :
                CabRequestsBelow(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
                                    (DirnBehaviourPair){D_Stop, EB_Idle}     ;
    case D_Down:
        return  CabRequestsBelow(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
                CabRequestsHere(e)  ? (DirnBehaviourPair){D_Up,   EB_DoorOpen} :
                CabRequestsAbove(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
                                    (DirnBehaviourPair){D_Stop, EB_Idle}     ;
    case D_Stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
        return  CabRequestsHere(e)  ? (DirnBehaviourPair){D_Stop, EB_DoorOpen} :
                CabRequestsAbove(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
                CabRequestsBelow(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
                                    (DirnBehaviourPair){D_Stop, EB_Idle}     ;
    default:
        return (DirnBehaviourPair){D_Stop, EB_Idle};
    }
}*/

func RequestsShouldStop(e Elevator) bool {
	switch e.Dirn {
	case elevio.MD_Down:
		return HallDownRequestHere(e) || RequestCabHere(e) || !CabRequestsBelow(e) //bool eller int?
	case elevio.MD_Up:
		return HallUpRequestHere(e) || RequestCabHere(e) || !CabRequestsAbove(e)
	case elevio.MD_Stop:
		fallthrough
	default:
		return true
	}
}

//i C
/*
int RequestsShouldStop(Elevator e){
    switch(e.dirn){
    case D_Down:
        return
            e.requests[e.floor][B_HallDown] ||
            e.requests[e.floor][B_Cab]      ||
            !CabRequestsBelow(e);
    case D_Up:
        return
            e.requests[e.floor][B_HallUp]   ||
            e.requests[e.floor][B_Cab]      ||
            !CabRequestsAbove(e);
    case D_Stop:
    default:
        return 1;
    }
}*/

func RequestsShouldClearImmediately(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
	return e.Floor == btn_floor && ((e.Dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) || (e.Dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) || (e.Dirn == elevio.MD_Stop || btn_type == elevio.BT_Cab))
}

//i C
/*
int RequestsShouldClearImmediately(Elevator e, int btn_floor, Button btn_type){
    switch(e.config.clearRequestVariant){
    case CV_All:
        return e.floor == btn_floor;
    case CV_InDirn:
        return
            e.floor == btn_floor &&
            (
                (e.dirn == D_Up   && btn_type == B_HallUp)    ||
                (e.dirn == D_Down && btn_type == B_HallDown)  ||
                e.dirn == D_Stop ||
                btn_type == B_Cab
            );
    default:
        return 0;
    }
}*/

func RequestsClearAtCurrentFloor(e Elevator) Elevator {
	/*
		    if e.Dirn == elevio.MD_Up {
				hallRequestCleared <- [2]int{e.Floor, 0}
			} else {
				hallRequestCleared <- [2]int{e.Floor, 1}
			}

			for btn := 0; btn < 3; btn++ {
				e.Requests[e.Floor][btn] = false
			}
	*/
	///*case CV_InDirn: sara
	e.Requests[e.Floor][elevio.BT_Cab] = false
	switch e.Dirn {
	case elevio.MD_Up:
		if !CabRequestsAbove(e) && !HallUpRequestHere(e) {
		 	e.Requests[e.Floor][elevio.BT_HallDown] = false
		 }
		// /*
		// 	if CabRequestsAbove(e) && e.Requests[e.Floor][elevio.BT_HallUp] && e.Requests[e.Floor][elevio.BT_HallDown] && !CabRequestsBelow(e){
		// 		e.Requests[e.Floor][elevio.BT_HallDown]=false
		// 		e.Behaviour=elevator.EB_DoorOpen

		// 	}*/
		// e.Requests[e.Floor][elevio.BT_HallUp] = false
		e.Requests[e.Floor][elevio.BT_HallUp] = false
		//if !RequestsAbove(e) {
		//	e.Requests[e.Floor][elevio.BT_HallDown] = false
		//}
	case elevio.MD_Down:
		 if !CabRequestsBelow(e) && !HallDownRequestHere(e) {
		 	e.Requests[e.Floor][elevio.BT_HallUp] = false
		 }
		e.Requests[e.Floor][elevio.BT_HallDown] = false
		//if !RequestsBelow(e) {
		//	e.Requests[e.Floor][elevio.BT_HallUp] = false
		//}
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
	fmt.Println("buttonToclear is:", buttonToclear)
	return buttonToclear

}

/*
Elevator RequestsClearAtCurrentFloor(Elevator e){

    switch(e.config.clearRequestVariant){
    case CV_All:
        for(Button btn = 0; btn < N_BUTTONS; btn++){
            e.requests[e.floor][btn] = 0;
        }
        break;

    case CV_InDirn:
        e.requests[e.floor][B_Cab] = 0;
        switch(e.dirn){
        case D_Up:
            if(!RequestsAbove(e) && !e.requests[e.floor][B_HallUp]){
                e.requests[e.floor][B_HallDown] = 0;
            }
            e.requests[e.floor][B_HallUp] = 0;
            break;

        case D_Down:
            if(!RequestsBelow(e) && !e.requests[e.floor][B_HallDown]){
                e.requests[e.floor][B_HallUp] = 0;
            }
            e.requests[e.floor][B_HallDown] = 0;
            break;

        case D_Stop:
        default:
            e.requests[e.floor][B_HallUp] = 0;
            e.requests[e.floor][B_HallDown] = 0;
            break;
        }
        break;

    default:
        break;
    }

    return e;
}*/
