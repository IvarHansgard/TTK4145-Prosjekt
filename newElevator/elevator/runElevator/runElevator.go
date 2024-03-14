package runElevator

import (
	"elevatorlib/elevator"
	"elevatorlib/elevator/requests"
	"elevatorlib/elevio"
	"elevatorlib/requestAsigner"
	"fmt"
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

/*func onDoorTimeout(e elevator.Elevator) elevator.Elevator {
	fmt.Println("Door timeout started")
	switch e.Behaviour {
	case elevator.EB_DoorOpen:
		pair := requests.RequestsChooseDirection(e)
		e.Dirn = pair.Dirn
		e.Behaviour = pair.Behaviour
		//elevio.SetDoorOpenLamp(false)
		fmt.Println("Door closed")
		fmt.Println(e)

		switch e.Behaviour {
		case elevator.EB_DoorOpen:
			doorTimeoutSignal.Reset(3 * time.Second)
			fmt.Println("Door Timout reset")

			e.Behaviour = elevator.EB_DoorOpen // door open
		case elevator.EB_Moving:
		case elevator.EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(e.Dirn)
			fmt.Println("Going to direction: ", e.Dirn)
			fmt.Println(e)
		}

	default:
		break
	}
	return e
}*/

/*
	func floorTimout(e elevator.Elevator, timeoutSignal chan bool, timer int) {
		timeout := timer
		fmt.Println("Floor timeout started")
		for {
			fmt.Println(timeout)
			select {
			case ts := <-timeoutSignal:
				if !ts {
					timeout = timer
				}
			default:
				if timeout == 0 {
					for i := 0; i < 4; i++ {
						for j := 0; j < 3; j++ {
							if e.Requests[i][j] {
								timeout = timer
							}
						}
					}
					if timeout == 0 {
						fmt.Println("Floor timeout")
						timeoutSignal <- true
					}
				}
				time.Sleep(1 * time.Second)
				timeout--
			}
		}
	}
*/
func onInitBetweenFloors(e elevator.Elevator) elevator.Elevator { //fikse denne

	fmt.Println("Elevator initialized between floors")

	e.Dirn = requests.RequestsChooseDirection(e).Dirn
	e.Behaviour = requests.RequestsChooseDirection(e).Behaviour
	// if e.Floor != 0 && e.Floor != 3 {
	// 	elevio.SetMotorDirection(elevio.MD_Down)
	// 	e.Dirn = elevio.MD_Down
	// 	e.Behaviour = elevator.EB_Moving
	// }
	return e
}

func localElevatorInit(id, port int) elevator.Elevator {
	fmt.Println("Initializing elevator ", id)

	elevio.Init("localhost:"+fmt.Sprint(port), 4)
	//elevio.Init("localhost:15657", 4)
	elevator := elevator.Elevator_init(id)
	setAllLights(elevator)
	elevio.SetDoorOpenLamp(false)

	return elevator
}

func RunLocalElevator(elevatorTx chan elevator.Elevator,
	newHallRequest chan elevio.ButtonEvent, assignedHallRequestsRx chan requestAsigner.HallRequests, chClearedHallRequests chan elevio.ButtonEvent, id, port int) {
	fmt.Println("Starting localElevator")

	localElevator := localElevatorInit(id, port)

	chButtonEvent := make(chan elevio.ButtonEvent)
	chFloor := make(chan int)
	chObstructionSwitch := make(chan bool)
	chStopButton := make(chan bool)

	doorTimeoutSignal := time.NewTimer(3 * time.Second)
	//chFloorTimeout := make(chan bool)

	go elevio.PollButtons(chButtonEvent)
	go elevio.PollFloorSensor(chFloor)
	go elevio.PollObstructionSwitch(chObstructionSwitch)
	go elevio.PollStopButton(chStopButton)
	//go floorTimout(localElevator, chFloorTimeout, 10)

	/*localElevator.Dirn = elevio.MD_Up
	localElevator.Behaviour = elevator.EB_Moving
	elevio.SetMotorDirection(localElevator.Dirn)*/
	if localElevator.Floor == -1 {
		localElevator = onInitBetweenFloors(localElevator)
	}
	for {
		select {
		/*
			case floorTimeout := <-chFloorTimeout:
					if floorTimeout {
						if localElevator.Floor != 0 {
							localElevator.Dirn = elevio.MotorDirection(elevio.MD_Down)
							elevio.SetMotorDirection(localElevator.Dirn)
						}
					}
					chFloorTimeout <- false
		*/
		case <-doorTimeoutSignal.C:

			if elevio.GetObstruction() {
				doorTimeoutSignal.Reset(3 * time.Second)
				break
			}

			fmt.Println("Door timed out")
			setAllLights(localElevator)
			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
				chClearedHallRequests <- requests.RequestClearHallRequestsAtCurrentFloor(localElevator)

				pair := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour
				//elevio.SetDoorOpenLamp(false)

				fmt.Println(localElevator)
				elevatorTx <- localElevator

				switch localElevator.Behaviour {
				case elevator.EB_DoorOpen:
					doorTimeoutSignal.Reset(3 * time.Second)
					fmt.Println("Door Timout reset")
					if localElevator.Requests[localElevator.Floor][0] || localElevator.Requests[localElevator.Floor][1] {
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
						chClearedHallRequests <- requests.RequestClearHallRequestsAtCurrentFloor(localElevator)
					} else {
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
					}
					setAllLights(localElevator)
					localElevator.Behaviour = elevator.EB_Idle
					//localElevator.Behaviour = elevator.EB_DoorOpen // door open
				case elevator.EB_Moving:
				case elevator.EB_Idle:
					elevio.SetDoorOpenLamp(false)
					fmt.Println("Door closed")
					if localElevator.Requests[localElevator.Floor][0] || localElevator.Requests[localElevator.Floor][1] {
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
						chClearedHallRequests <- requests.RequestClearHallRequestsAtCurrentFloor(localElevator)

					} else {
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
					}
					elevio.SetMotorDirection(localElevator.Dirn)
					fmt.Println("Going to direction: ", localElevator.Dirn)
					fmt.Println(localElevator)
				}

			default:
				break
			}
			elevatorTx <- localElevator

		case HallRequests := <-assignedHallRequestsRx:
			switch id {

			case 0:
				for i := 0; i < 4; i++ {
					for j := 0; j < 2; j++ {
						localElevator.Requests[i][j] = HallRequests["one"][i][j]

					}
				}

				fmt.Println("Assigned hall requests: ", HallRequests["one"])
				fmt.Println("All requests: ", localElevator.Requests)
				dirn := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = dirn.Dirn
				localElevator.Behaviour = dirn.Behaviour
				setAllLights(localElevator)
				fmt.Println("set direction", localElevator.Dirn)
				elevio.SetMotorDirection(localElevator.Dirn)

			case 1:
				for i := 0; i < 4; i++ {
					for j := 0; j < 2; j++ {
						localElevator.Requests[i][j] = HallRequests["two"][i][j]

					}
				}
				fmt.Println("Assigned hall requests: ", HallRequests["two"])
				dirn := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = dirn.Dirn
				localElevator.Behaviour = dirn.Behaviour
				setAllLights(localElevator)
				elevio.SetMotorDirection(localElevator.Dirn)

			case 2:
				for i := 0; i < 4; i++ {
					for j := 0; j < 2; j++ {
						localElevator.Requests[i][j] = HallRequests["three"][i][j]

					}
				}
				fmt.Println("Assigned hall requests: ", HallRequests["three"])
				dirn := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = dirn.Dirn
				localElevator.Behaviour = dirn.Behaviour
				setAllLights(localElevator)
				elevio.SetMotorDirection(localElevator.Dirn)
			}
		case button := <-chButtonEvent:

			fmt.Println("Button press detected")
			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				fmt.Println("Door is open")
				if requests.RequestsShouldClearImmediately(localElevator, button.Floor, button.Button) {
					fmt.Println("Should clear immediately")
					doorTimeoutSignal.Reset(3 * time.Second)
					localElevator.Behaviour = elevator.EB_DoorOpen
					elevatorTx <- localElevator

				} else {
					if button.Button == elevio.BT_HallDown || button.Button == elevio.BT_HallUp {
						fmt.Println("new hall request")
						newHallRequest <- button
					} else {
						fmt.Println("new cab request")
						localElevator.Requests[button.Floor][button.Button] = true
						elevatorTx <- localElevator
					}
				}

			case elevator.EB_Moving:
				fmt.Println("Elevator is moving.")
				//elevio.SetMotorDirection(localElevator.Dirn)
				if button.Button == elevio.BT_HallDown || button.Button == elevio.BT_HallUp {
					newHallRequest <- button
				} else {
					localElevator.Requests[button.Floor][button.Button] = true
					elevatorTx <- localElevator
				}

			case elevator.EB_Idle:
				fmt.Println("Elevator is idle.")
				if button.Button == elevio.BT_HallDown || button.Button == elevio.BT_HallUp {
					fmt.Println("Sent hall request:", button)
					newHallRequest <- button

				} else {
					localElevator.Requests[button.Floor][button.Button] = true
					elevatorTx <- localElevator
				}
				pair := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour
				elevatorTx <- localElevator
				switch pair.Behaviour {
				case elevator.EB_DoorOpen:
					fmt.Println("Door is open")
					elevio.SetDoorOpenLamp(true)
					doorTimeoutSignal.Reset(3 * time.Second)
					/*if localElevator.Requests[localElevator.Floor][0] || localElevator.Requests[localElevator.Floor][1] {
						chClearedHallRequests <- requests.RequestClearHallRequestsAtCurrentFloor(localElevator)
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
					} else {
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
					}*/
					elevatorTx <- localElevator
				case elevator.EB_Moving:
					fmt.Println("Elevator is moving.")

					elevio.SetMotorDirection(localElevator.Dirn)
					elevatorTx <- localElevator
				case elevator.EB_Idle:
					break
				}
			}
			setAllLights(localElevator)
			elevatorTx <- localElevator //check hall request ikke oppdater nå hall

		case floor := <-chFloor:
			localElevator.Floor = floor

			fmt.Println("elevator is on floor: ", floor)
			elevio.SetFloorIndicator(localElevator.Floor)
			elevatorTx <- localElevator

			switch localElevator.Behaviour {
			case elevator.EB_Moving:
				if requests.RequestsShouldStop(localElevator) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					if localElevator.Requests[localElevator.Floor][0] || localElevator.Requests[localElevator.Floor][1] {
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
						chClearedHallRequests <- requests.RequestClearHallRequestsAtCurrentFloor(localElevator) //flyttet fra 360

					} else {
						localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
					}
					fmt.Println("Door open")
					doorTimeoutSignal.Reset(3 * time.Second)
					setAllLights(localElevator)
					localElevator.Behaviour = elevator.EB_DoorOpen
					//requests.RequestsClearAtCurrentFloor(localElevator)
					elevatorTx <- localElevator
					//buttonevent := requests.RequestClearHallRequestsAtCurrentFloor(localElevator)
					//fmt.Println("clearing request: ", buttonevent)
					//chClearedHallRequests <- requests.RequestClearHallRequestsAtCurrentFloor(localElevator)

					fmt.Println("dirn is: ", localElevator.Dirn)
				}
				break
			case elevator.EB_Idle:
				for i := 0; i < 4; i++ {
					for j := 0; j < 3; j++ {
						if localElevator.Requests[i][j] {
							localElevator.Dirn = requests.RequestsChooseDirection(localElevator).Dirn
							localElevator.Behaviour = requests.RequestsChooseDirection(localElevator).Behaviour
							elevio.SetMotorDirection(localElevator.Dirn)
							break
						}
					}
				}
				elevatorTx <- localElevator

			default:
				break
			}

		case <-chStopButton:
			//set behavior
			fmt.Printf("%+v\n", localElevator)
			for f := 0; f < 4; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
					elevio.SetMotorDirection(elevio.MD_Stop)
				}
			}
		}
	}
}
