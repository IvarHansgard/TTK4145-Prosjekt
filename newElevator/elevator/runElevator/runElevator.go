package runElevator

import (
	"elevatorlib/elevator"
	"elevatorlib/elevator/requests"
	"elevatorlib/elevio"
	"elevatorlib/requestAsigner"
	"time"
)

func setAllLights(e elevator.Elevator) {
	for f := 0; f < 4; f++ {
		for btn := 0; btn < 3; btn++ {
			switch btn {
			case 0:
				b := elevio.BT_HallUp
				elevio.SetButtonLamp(b, f, e.Requests[f][btn])
			case 1:
				b := elevio.BT_HallDown
				elevio.SetButtonLamp(b, f, e.Requests[f][btn])
			case 2:
				b := elevio.BT_Cab
				elevio.SetButtonLamp(b, f, e.Requests[f][btn])

			}
		}
	}
}

func onDoorTimeout(e elevator.Elevator) {
	switch e.Behaviour {
	case "EB_DoorOpen":
		pair := requests.RequestsChooseDirection(e)
		e.Dirn = pair.Dirn
		e.Behaviour = pair.Behaviour

		switch e.Behaviour {
		case "EB_DoorOpen":
			time.Sleep(3 * time.Second)
			e.Requests = requests.RequestsClearAtCurrentFloor(e)
			setAllLights(e)
			break

		case "EB_Moving":

		case "EB_Idle":
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(e.dirn)
			break
		}
		break

	default:
		break
	}

}

func onInitBetweenFloors(e elevator.Elevator) {
	e.Dirn = requests.RequestsChooseDirection(e)
	e.Behaviour = requests.RequestsChooseDirection(e)
	elevio.SetMotorDirection(e.Dirn)
}

func RunLocalElevator(chActiveElevators chan []elevator.Elevator, elevatorTx chan elevator.Elevator, hallRequestsRx chan requestAsigner.HallRequests, id int) {
	elevio.Init("localhost:15657", 4)

	localElevator := elevator.Elevator_init(id)

	chButtonEvent := make(chan elevio.ButtonEvent)
	chFloor := make(chan int)
	chObstructionSwitch := make(chan bool)
	chStopButton := make(chan bool)

	go elevio.PollButtons(chButtonEvent)
	go elevio.PollFloorSensor(chFloor)
	go elevio.PollObstructionSwitch(chObstructionSwitch)
	go elevio.PollStopButton(chStopButton)

	if -1 == <-chFloor {
		onInitBetweenFloors(localElevator)
	}

	for {
		select {
		case HallRequests := <-hallRequestsRx:
			switch id {
			case 0:
				for i := 0; i < 4; i++ {
					for j := 0; j < 2; j++ {
						localElevator.Requests[i][j] = HallRequests["one"][i][j]
					}
				}
			case 1:
				for i := 0; i < 4; i++ {
					for j := 0; j < 2; j++ {
						localElevator.Requests[i][j] = HallRequests["two"][i][j]
					}
				}
			case 2:
				for i := 0; i < 4; i++ {
					for j := 0; j < 2; j++ {
						localElevator.Requests[i][j] = HallRequests["three"][i][j]
					}
				}
			}

		case button := <-chButtonEvent:
			localElevator.Requests[button.Floor][button.Button] = true
			elevatorTx <- localElevator //check hall request ikke oppdater nå hall

		case floor := <-chFloor:
			localElevator.Floor = floor
			elevio.SetFloorIndicator(localElevator.Floor)

			switch localElevator.Behaviour {
			case "EB_Moving":
				if requests.RequestsShouldStop(localElevator) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					requests.RequestsClearAtCurrentFloor(localElevator)

					time.Sleep(3 * time.Second)
					setAllLights(localElevator)
					localElevator.Behaviour = "EB_DoorOpen"

					elevatorTx <- localElevator

					time.Sleep(10 * time.Second)
					onDoorTimeout(localElevator)
				}
				elevatorTx <- localElevator

				break
			default:
				break
			}

		case obstruction := <-chObstructionSwitch:
			if obstruction {
				elevio.SetMotorDirection(elevio.MD_Stop)
				//set behavior

			} else {
				elevio.SetMotorDirection(localElevator.Dirn)
				//set behavior

			}
		case stop := <-chStopButton:
			//set behavior
			for f := 0; f < 4; f++ {
				for btn := elevio.ButtonType(0); btn < 3; btn++ {
					elevio.SetButtonLamp(btn, f, false)
				}
			}
		}

	}
}
