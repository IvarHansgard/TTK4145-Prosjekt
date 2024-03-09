package requests

import (
	. "elevatorlib/elevator"
	"elevatorlib/elevio"
)

type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection
	Behaviour ElevatorBehaviour
}

func requests_above(e Elevator) int {
	for f := e.floor + 1; f < 4; f++ {
		for btn := 0; btn < 3; btn++ {
			if e.requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

// i C
/*
static int requests_above(Elevator e){
    for(int f = e.floor+1; f < N_FLOORS; f++){
        for(int btn = 0; btn < N_BUTTONS; btn++){
            if(e.requests[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}*/

func Requests_below(e Elevator) int {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.request[f][btn] {
				return 1
			}
		}
	}
	return 0
}

// i C
/*
static int requests_below(Elevator e){
    for(int f = 0; f < e.floor; f++){
        for(int btn = 0; btn < N_BUTTONS; btn++){
            if(e.requests[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}*/

func Requests_here(e Elevator) int {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.requests[e.floor][btn] {
			return 1
		}
	}
	return 0
}

/*
// i C:
static int requests_here(Elevator e){
    for(int btn = 0; btn < N_BUTTONS; btn++){
        if(e.requests[e.floor][btn]){
            return 1;
        }
    }
    return 0;
}*/

func Requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case D_up:
		if requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if request_here(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else {
			DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Down:
		if requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if request_here(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else if requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else {
			DirnBehaviourPair{D_Stop, EB_Idle}
		}

	case D_Stop:
		if request_here(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requests_below(e) {
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
DirnBehaviourPair requests_chooseDirection(Elevator e){
    switch(e.dirn){
    case D_Up:
        return  requests_above(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
                requests_here(e)  ? (DirnBehaviourPair){D_Down, EB_DoorOpen} :
                requests_below(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
                                    (DirnBehaviourPair){D_Stop, EB_Idle}     ;
    case D_Down:
        return  requests_below(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
                requests_here(e)  ? (DirnBehaviourPair){D_Up,   EB_DoorOpen} :
                requests_above(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
                                    (DirnBehaviourPair){D_Stop, EB_Idle}     ;
    case D_Stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
        return  requests_here(e)  ? (DirnBehaviourPair){D_Stop, EB_DoorOpen} :
                requests_above(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
                requests_below(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
                                    (DirnBehaviourPair){D_Stop, EB_Idle}     ;
    default:
        return (DirnBehaviourPair){D_Stop, EB_Idle};
    }
}*/

func Requests_shouldStop(e Elevator) int {
	switch e.dirn {
	case D_Down:
		return e.requests[e.floor][B_HallDown] || e.requests[e.floor][B_Cab] || !requests_below(e) //bool eller int?
	case D_Up:
		return e.requests[e.floor][B_HallUp] || e.requests[e.floor][B_Cab] || !requests_above(e)
	case D_Stop:
		fallthrough
	default:
		return 1
	}
}

//i C
/*
int requests_shouldStop(Elevator e){
    switch(e.dirn){
    case D_Down:
        return
            e.requests[e.floor][B_HallDown] ||
            e.requests[e.floor][B_Cab]      ||
            !requests_below(e);
    case D_Up:
        return
            e.requests[e.floor][B_HallUp]   ||
            e.requests[e.floor][B_Cab]      ||
            !requests_above(e);
    case D_Stop:
    default:
        return 1;
    }
}*/

func Requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type Button) {
	switch e.config.clearRequestVariant {
	case CV_All:
		return e.floor == btn_floor
	case CV_InDirn:
		return e.floor == btn_floor && ((e.dirn == D_Up && btn_type == B_HallUp) ||
			(e.dirn == D_Down && btn_type == B_HallDown) ||
			e.dirn == D_Stop ||
			btn_type == B_Cab)
	default:
		return 0
	}
}

//i C
/*
int requests_shouldClearImmediately(Elevator e, int btn_floor, Button btn_type){
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

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	switch e.Config.clearRequestVariant {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.Requests[e.floor][btn] = 0
		}
		break
	case CV_InDirn:
		e.Request[e.floor][B_Cab] = 0
		switch e.Dirn {
		case D_Up:
			if !requests_above(e) && !e.Request[e.floor][B_HallUp] {
				e.Request[e.floor][B_HallDown] = 0
			}
			e.Request[e.floor][B_HallUp] = 0
		case D_Down:
			if !requests_below(e) && !e.Requests[e.floor][B_HallDown] {
				e.requests[e.floor][B_HallUp] = 0
			}
			e.Requests[e.floor][B_HallDown] = 0
		case D_Stop:
			fallthrough
		default:
			e.Requests[e.floor][B_HallUp] = 0
			e.Requests[e.floor][B_HallDown] = 0
		}
	default:
		break
	}
	return e

}

/*
Elevator requests_clearAtCurrentFloor(Elevator e){

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
            if(!requests_above(e) && !e.requests[e.floor][B_HallUp]){
                e.requests[e.floor][B_HallDown] = 0;
            }
            e.requests[e.floor][B_HallUp] = 0;
            break;

        case D_Down:
            if(!requests_below(e) && !e.requests[e.floor][B_HallDown]){
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
