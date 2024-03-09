package requests

import (
	. "elevatorlib/elevator"
	"elevatorlib/elevio"
)

type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection
	Behaviour ElevatorBehaviour
}

func RequestsAbove(e Elevator) int {
	for f := e.Floor + 1; f < 4; f++ {
		for btn := 0; btn < 3; btn++ {
			if e.Requests[f][btn] {
				return 1
			}
		}
	}
	return 0
} 

// i C
/*
static int RequestsAbove(Elevator e){
    for(int f = e.floor+1; f < N_FLOORS; f++){
        for(int btn = 0; btn < N_BUTTONS; btn++){
            if(e.requests[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}*/

func RequestsBelow(e Elevator) int {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

// i C
/*
static int RequestsBelow(Elevator e){
    for(int f = 0; f < e.floor; f++){
        for(int btn = 0; btn < N_BUTTONS; btn++){
            if(e.requests[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}*/

func RequestsHere(e Elevator) int {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return 1
		}
	}
	return 0
}

/*
// i C:
static int RequestsHere(Elevator e){
    for(int btn = 0; btn < N_BUTTONS; btn++){
        if(e.requests[e.floor][btn]){
            return 1;
        }
    }
    return 0;
}*/

func RequestsChooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case D_up:
		if RequestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else {
			DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Down:
		if RequestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else {
			DirnBehaviourPair{D_Stop, EB_Idle}
		}

	case D_Stop:
		if RequestsHere(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			DirnBehaviourPair{D_Stop, EB_Idle}
		}

	default:
		return DirnBehaviourPair{D_Stop, EB_Idle}
	}
}

//i C
/*
DirnBehaviourPair RequestsChooseDirection(Elevator e){
    switch(e.dirn){
    case D_Up:
        return  RequestsAbove(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
                RequestsHere(e)  ? (DirnBehaviourPair){D_Down, EB_DoorOpen} :
                RequestsBelow(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
                                    (DirnBehaviourPair){D_Stop, EB_Idle}     ;
    case D_Down:
        return  RequestsBelow(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
                RequestsHere(e)  ? (DirnBehaviourPair){D_Up,   EB_DoorOpen} :
                RequestsAbove(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
                                    (DirnBehaviourPair){D_Stop, EB_Idle}     ;
    case D_Stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
        return  RequestsHere(e)  ? (DirnBehaviourPair){D_Stop, EB_DoorOpen} :
                RequestsAbove(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
                RequestsBelow(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
                                    (DirnBehaviourPair){D_Stop, EB_Idle}     ;
    default:
        return (DirnBehaviourPair){D_Stop, EB_Idle};
    }
}*/

func RequestsShouldStop(e Elevator) int {
	switch e.Dirn {
	case D_Down:
		return e.Requests[e.Floor][B_HallDown] || e.Requests[e.Floor][B_Cab] || !RequestsBelow(e) //bool eller int?
	case D_Up:
		return e.Requests[e.Floor][B_HallUp] || e.Requests[e.Floor][B_Cab] || !RequestsAbove(e)
	case D_Stop:
		fallthrough
	default:
		return 1
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
            !RequestsBelow(e);
    case D_Up:
        return
            e.requests[e.floor][B_HallUp]   ||
            e.requests[e.floor][B_Cab]      ||
            !RequestsAbove(e);
    case D_Stop:
    default:
        return 1;
    }
}*/

func RequestsShouldClearImmediately(e Elevator, btn_floor int, btn_type Button) {
	switch e.config.clearRequestVariant {
	case CV_All:
		return e.Floor == btn_floor
	case CV_InDirn:
		return e.Floor == btn_floor && ((e.Dirn == D_Up && btn_type == B_HallUp) ||
			(e.Dirn == D_Down && btn_type == B_HallDown) ||
			e.Dirn == D_Stop ||
			btn_type == B_Cab)
	default:
		return 0
	}
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
	switch e.Config.clearRequestVariant {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.Requests[e.Floor][btn] = 0
		}
		break
	case CV_InDirn:
		e.Requests[e.Floor][B_Cab] = 0
		switch e.Dirn {
		case D_Up:
			if !RequestsAbove(e) && !e.Requests[e.Floor][B_HallUp] {
				e.Requests[e.Floor][B_HallDown] = 0
			}
			e.Requests[e.Floor][B_HallUp] = 0
		case D_Down:
			if !RequestsBelow(e) && !e.Requests[e.Floor][B_HallDown] {
				e.Requests[e.Floor][B_HallUp] = 0
			}
			e.Requests[e.Floor][B_HallDown] = 0
		case D_Stop:
			fallthrough
		default:
			e.Requests[e.Floor][B_HallUp] = 0
			e.Requests[e.Floor][B_HallDown] = 0
		}
	default:
		break
	}
	return e

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
