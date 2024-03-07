package runElevator

import (
	elevator "ElevatorLib/Elevator"
	"ElevatorLib/elevator"
	"ElevatorLib/elevio"
	"ElevatorLib/requests"
	"time"
)

func setAllLights(e chan elevator.Elevator) {
	for f := 0; f < 4; f++ {
		for btn := 0; btn < 3; btn++ {
			elevio.SetButtonLamp(btn, f, e.Requests[f][btn])
		}
	}
}

func onDoorTimeout(e chan elevator.Elevator) {
	switch e.Behaviour {
	case EB_DoorOpen:
		pair := requests.Requests_chooseDirection(e)
		e.Dirn = pair.Dirn
		e.Behaviour = pair.Behaviour

		switch e.behaviour {
		case EB_DoorOpen:
			time.Sleep(3 * time.Second)
			e.Requests = request.requests_clearAtCurrentFloor(e)
			elevio.setAllLights(e)
			break

		case EB_Moving:

		case EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(e.dirn)
			break
		}
		break

	default:
		break
	}

}

func onInitBetweenFloors(e chan elevator.Elevator) {
	e.Dirn = requests.Requests_chooseDirection(e).Dirn
	e.Behaviour = requests.Requests_chooseDirection(e).Behaviour
	elevio.SetMotorDirection(e.Dirn)
}

func RunElevator(e chan elevator.Elevator) {
	elvio.Init("localhost:15657", 4)

	chButtonEvent := make(chan elevio.ButtonEvent)
	chFloor := make(chan int)
	chObstructionSwitch := make(chan bool)
	chStopButton := make(chan bool)

	go elevio.PollButtons(buttonEvent)
	go elevio.PollFloorSensor(floor)
	go elevio.PollObstructionSwitch(chObstructionSwitch)
	go elevio.PollStopButton(chStopButton)

	if floor == -1 {
		onInitBetweenFloors(e)
	}

	for {
		select {
		case button := <-chButtonEvent:
			e.Requests[button.Floor][button.ButtonType] = true
			elevatorTx <- e

		case floor := <-chFloor:
			e.Floor = floor
			elevio.SetFloorIndicator(e.Floor)

			switch e.Behaviour {
			case EB_Moving:
				if requests.Requests_shouldStop(e) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					requests.Requests_clearAtCurrentFloor(e)

					elevatorTx <- e

					time.Sleep(3 * time.second)
					elevio.setAllLights(e)
					e.behaviour = EB_DoorOpen
					time.Sleep(10 * time.Second())
					onDoorTimeout(e)
				}
				break
			default:
				break
			}

		case obstruction := <-chobstructionSwitch:
			if obstruction {
				elevio.SetMotorDirection(elevio.MD_Stop)

			} else {
				elevio.SetMotorDirection(e.Dirn)
			}
		case stop := <-chStopButton:
			for f := 0; f < 4; f++ {
				for btn := elevio.ButtonType(0); btn < 3; btn++ {
					elvio.SetButtonLamp(btn, f, false)
				}
			}
		}

	}
}
