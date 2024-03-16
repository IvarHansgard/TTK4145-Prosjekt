package runElevator

import (
	"elevatorlib/elevator"
	"elevatorlib/elevator/requests"
	"elevatorlib/elevio"
	"elevatorlib/requestAsigner"
	"fmt"
	"strconv"
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

func localElevatorInit(id, port int, strId string) elevator.Elevator {
	fmt.Println("Initializing elevator ", id)

	elevio.Init("localhost:"+fmt.Sprint(port), 4)
	//elevio.Init("localhost:15657", 4)
	elevator := elevator.Elevator_init(id, strId)
	setAllLights(elevator)
	elevio.SetDoorOpenLamp(false)

	return elevator
}

func RunLocalElevator(chElevatorTx chan elevator.Elevator, chNewHallRequestTx chan elevio.ButtonEvent, chAssignedHallRequestsRx chan requestAsigner.HallRequests,
	chHallRequestClearedTx chan elevio.ButtonEvent, strId string, port int, chStopButtonPressed chan bool, chSetButtonLightRx chan elevio.ButtonEvent, chSetButtonLightTx chan elevio.ButtonEvent) {
	fmt.Println("Starting localElevator")
	id, err := strconv.Atoi(strId)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	localElevator := localElevatorInit(id, port, strId)

	chButtonEvent := make(chan elevio.ButtonEvent)
	chFloor := make(chan int)
	chObstructionSwitch := make(chan bool)
	chStopButton := make(chan bool)

	doorTimeoutSignal := time.NewTimer(3 * time.Second)

	go elevio.PollButtons(chButtonEvent)
	go elevio.PollFloorSensor(chFloor)
	go elevio.PollObstructionSwitch(chObstructionSwitch)
	go elevio.PollStopButton(chStopButton)

	if elevio.GetFloor() == -1 {
		fmt.Println("elevator initialized on ", elevio.GetFloor())
		elevio.SetMotorDirection(elevio.MD_Down)
		localElevator.Dirn = elevio.MD_Down
		localElevator.Behaviour = elevator.EB_Moving
		chElevatorTx <- localElevator
	}
	for {
		select {
		case <-doorTimeoutSignal.C:
			if elevio.GetObstruction() {
				doorTimeoutSignal.Reset(3 * time.Second)
				break
			}

			fmt.Println("Door timed out")
			setAllLights(localElevator)

			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:

				pair := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour

				switch localElevator.Behaviour {
				case elevator.EB_DoorOpen:
					doorTimeoutSignal.Reset(3 * time.Second)
					if localElevator.Requests[localElevator.Floor][0] || localElevator.Requests[localElevator.Floor][1] {
						if localElevator.Requests[localElevator.Floor][0] && localElevator.Requests[localElevator.Floor][1] {
							if localElevator.Dirn == elevio.MD_Up {
								localElevator.Dirn = elevio.MD_Down
								chHallRequestClearedTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
								localElevator.Requests[localElevator.Floor][requests.HallRequestsClearAtCurrentFloor(localElevator).Button] = false
								chSetButtonLightTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
								localElevator.Dirn = elevio.MD_Up
							} else if localElevator.Dirn == elevio.MD_Down {
								localElevator.Dirn = elevio.MD_Up
								chHallRequestClearedTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
								localElevator.Requests[localElevator.Floor][requests.HallRequestsClearAtCurrentFloor(localElevator).Button] = false
								chSetButtonLightTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
								localElevator.Dirn = elevio.MD_Down
							}
						}
						localElevator.Requests[localElevator.Floor][requests.HallRequestsClearAtCurrentFloor(localElevator).Button] = false
						chHallRequestClearedTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
					}
					localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
					setAllLights(localElevator)
					break
				case elevator.EB_Moving:
					fallthrough
				case elevator.EB_Idle:
					elevio.SetDoorOpenLamp(false)
					fmt.Println("Door closed")
					if localElevator.Requests[localElevator.Floor][0] || localElevator.Requests[localElevator.Floor][1] {
						if localElevator.Requests[localElevator.Floor][0] && localElevator.Requests[localElevator.Floor][1] {
							if localElevator.Dirn == elevio.MD_Up {
								localElevator.Dirn = elevio.MD_Down
								chHallRequestClearedTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
								localElevator.Requests[localElevator.Floor][requests.HallRequestsClearAtCurrentFloor(localElevator).Button] = false
								chSetButtonLightTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
								localElevator.Dirn = elevio.MD_Up
							} else if localElevator.Dirn == elevio.MD_Down {
								localElevator.Dirn = elevio.MD_Up
								chHallRequestClearedTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)

								localElevator.Requests[localElevator.Floor][requests.HallRequestsClearAtCurrentFloor(localElevator).Button] = false
								chSetButtonLightTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
								localElevator.Dirn = elevio.MD_Down
							}
						}
						time.Sleep(50 * time.Millisecond)
						chHallRequestClearedTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
					}
					localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
					elevio.SetMotorDirection(localElevator.Dirn)
					break
				}
				break
			default:
				break
			}
			chElevatorTx <- localElevator

		case HallRequests := <-chAssignedHallRequestsRx:
			for i := 0; i < 4; i++ {
				for j := 0; j < 2; j++ {
					localElevator.Requests[i][j] = HallRequests[strId][i][j]
				}
			}
			dirn := requests.RequestsChooseDirection(localElevator)
			localElevator.Dirn = dirn.Dirn
			localElevator.Behaviour = dirn.Behaviour
			if localElevator.Dirn == elevio.MD_Stop {
				if requests.HallUpRequestHere(localElevator) || requests.HallDownRequestHere(localElevator) {
					fmt.Println("door open")
					localElevator.Behaviour = elevator.EB_DoorOpen
					elevio.SetDoorOpenLamp(true)
					doorTimeoutSignal.Reset(3 * time.Second)
				}
			} else {
				elevio.SetMotorDirection(localElevator.Dirn)
			}
			setAllLights(localElevator)
			//set lys for hall requests
			for i := 0; i < len(HallRequests); i++ {
				for j := 0; j < 4; j++ {
					for k := 0; k < 2; k++ {
						if HallRequests[strconv.Itoa(i)][j][k] {
							button := elevio.ButtonType(k)
							elevio.SetButtonLamp(button, j, true)
						}
					}
				}
			}
			chElevatorTx <- localElevator

		case temp := <-chSetButtonLightRx:
			//slå av lys for hall requests fullført av andre heiser
			elevio.SetButtonLamp(temp.Button, temp.Floor, false)

		case button := <-chButtonEvent:
			fmt.Println(localElevator)
			//fmt.Println("Button press detected")
			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				fmt.Println("Door is open")
				if requests.RequestsShouldClearImmediately(localElevator, button.Floor, button.Button) {
					doorTimeoutSignal.Reset(3 * time.Second)
				} else {
					if button.Button == elevio.BT_HallDown || button.Button == elevio.BT_HallUp {
						chNewHallRequestTx <- button
					} else {
						localElevator.Requests[button.Floor][button.Button] = true

					}
				}
				chElevatorTx <- localElevator

			case elevator.EB_Moving:
				//elevio.SetMotorDirection(localElevator.Dirn)
				if button.Button == elevio.BT_HallDown || button.Button == elevio.BT_HallUp {
					chNewHallRequestTx <- button
				} else {
					localElevator.Requests[button.Floor][button.Button] = true
					chElevatorTx <- localElevator
				}

			case elevator.EB_Idle:
				if button.Button == elevio.BT_HallDown || button.Button == elevio.BT_HallUp {
					chNewHallRequestTx <- button
				} else {
					localElevator.Requests[button.Floor][button.Button] = true
				}

				pair := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour

				chElevatorTx <- localElevator

				switch pair.Behaviour {
				case elevator.EB_DoorOpen:
					fmt.Println("Door is open")
					elevio.SetDoorOpenLamp(true)
					doorTimeoutSignal.Reset(3 * time.Second)
					localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
					chElevatorTx <- localElevator
					break
				case elevator.EB_Moving:
					elevio.SetMotorDirection(localElevator.Dirn)
					break
				case elevator.EB_Idle:
					break
				}
			}
			setAllLights(localElevator)
			chElevatorTx <- localElevator

		case floor := <-chFloor:
			localElevator.Floor = floor
			if localElevator.Requests[floor][0] {

			}
			if localElevator.Requests[floor][1] {

			}
			fmt.Println("elevator is on floor: ", floor)
			elevio.SetFloorIndicator(localElevator.Floor)
			switch localElevator.Behaviour {
			case elevator.EB_Moving:
				if requests.RequestsShouldStop(localElevator) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					if localElevator.Requests[localElevator.Floor][0] || localElevator.Requests[localElevator.Floor][1] {
						chHallRequestClearedTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
						chSetButtonLightTx <- requests.HallRequestsClearAtCurrentFloor(localElevator)
						localElevator.Requests[localElevator.Floor][requests.HallRequestsClearAtCurrentFloor(localElevator).Button] = false
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
						doorTimeoutSignal.Reset(3 * time.Second)
					} else {
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
						doorTimeoutSignal.Reset(3 * time.Second)
					}
					setAllLights(localElevator)
					localElevator.Behaviour = elevator.EB_DoorOpen
					chElevatorTx <- localElevator
				}
				break
			case elevator.EB_Idle:
			default:
				break
			}
			chElevatorTx <- localElevator

		case <-chStopButton:
			//set behavior
			fmt.Printf("%+v\n", localElevator)
			for f := 0; f < 4; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					localElevator.Requests[f][b] = false
					elevio.SetButtonLamp(b, f, false)
				}
			}
			localElevator.Dirn = elevio.MD_Stop
			localElevator.Behaviour = elevator.EB_Disconnected
			elevio.SetMotorDirection(localElevator.Dirn)
			elevio.SetDoorOpenLamp(false)
			chElevatorTx <- localElevator
			chStopButtonPressed <- true
		}
	}
}
