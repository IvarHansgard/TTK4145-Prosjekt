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

func onDoorTimeout(e elevator.Elevator, hallRequestCleared chan [2]int) elevator.Elevator {
	fmt.Println("Door timeout started")
	switch e.Behaviour {
	case elevator.EB_DoorOpen:
		pair := requests.RequestsChooseDirection(e)
		e.Dirn = pair.Dirn
		e.Behaviour = pair.Behaviour
		elevio.SetDoorOpenLamp(false)
		fmt.Println("Door closed")

		switch e.Behaviour {
		case elevator.EB_DoorOpen:
			e.Behaviour = elevator.EB_Idle
			fmt.Println("door timer restart")

			time.Sleep(3 * time.Second)
			e = requests.RequestsClearAtCurrentFloor(e, hallRequestCleared)
			fmt.Println(e)
			onDoorTimeout(e, hallRequestCleared)
			fmt.Println(e)

		case elevator.EB_Moving:
			break
		case elevator.EB_Idle:
			elevio.SetDoorOpenLamp(false)
			fmt.Println("Door closed")
			elevio.SetMotorDirection(e.Dirn)
			fmt.Println("Going to direction: ", e.Dirn)
			fmt.Println(e)
		}

	default:
		break
	}
	return e
}

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

	//e.Dirn = requests.RequestsChooseDirection(e).Dirn
	//e.Behaviour = requests.RequestsChooseDirection(e).Behaviour

	elevio.SetMotorDirection(elevio.MD_Down)
	e.Dirn = elevio.MD_Down
	e.Behaviour = elevator.EB_Moving

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

func RunLocalElevator(chActiveElevators chan []elevator.Elevator, elevatorTx chan elevator.Elevator,
	localHallRequestsTx chan [][2]bool, assignedHallRequestsRx chan requestAsigner.HallRequests, hallRequestCleared chan [2]int, id, port int) {
	fmt.Println("Starting localElevator")
	localElevator := localElevatorInit(id, port)
	localHallRequests := [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}

	chButtonEvent := make(chan elevio.ButtonEvent)
	chFloor := make(chan int)
	chObstructionSwitch := make(chan bool)
	chStopButton := make(chan bool)

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
		case HallRequests := <-assignedHallRequestsRx:
			switch id {
			case 0:
				for i := 0; i < 4; i++ {
					for j := 0; j < 2; j++ {
						localElevator.Requests[i][j] = HallRequests["one"][i][j]

					}
				}
				fmt.Println("Assigned hall requests: ", HallRequests["one"])
				dirn := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = dirn.Dirn
				localElevator.Behaviour = dirn.Behaviour
				setAllLights(localElevator)
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
					time.Sleep(3 * time.Second)
					localElevator = onDoorTimeout(localElevator, hallRequestCleared)
					elevatorTx <- localElevator

				} else {
					if button.Button == elevio.BT_HallDown || button.Button == elevio.BT_HallUp {
						localHallRequests[button.Floor][button.Button] = true
						localHallRequestsTx <- localHallRequests
					} else {
						localElevator.Requests[button.Floor][button.Button] = true
						elevatorTx <- localElevator
					}
				}

			case elevator.EB_Moving:
				fmt.Println("Elevator is moving.")
				if button.Button == elevio.BT_HallDown || button.Button == elevio.BT_HallUp {
					localHallRequests[button.Floor][button.Button] = true
					localHallRequestsTx <- localHallRequests
				} else {
					localElevator.Requests[button.Floor][button.Button] = true
					elevatorTx <- localElevator
				}

			case elevator.EB_Idle:
				fmt.Println("Elevator is idle.")
				if button.Button == elevio.BT_HallDown || button.Button == elevio.BT_HallUp {
					localHallRequests[button.Floor][button.Button] = true
					localHallRequestsTx <- localHallRequests
					elevatorTx <- localElevator

				} else {
					localElevator.Requests[button.Floor][button.Button] = true
					elevatorTx <- localElevator
				}
				pair := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour
				switch pair.Behaviour {
				case elevator.EB_DoorOpen:
					fmt.Println("Door is open")
					elevio.SetDoorOpenLamp(true)
					time.Sleep(3 * time.Second)
					localElevator = onDoorTimeout(localElevator, hallRequestCleared)
					localElevator = requests.RequestsClearAtCurrentFloor(localElevator, hallRequestCleared)
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
			elevio.SetFloorIndicator(localElevator.Floor)
			elevatorTx <- localElevator

			switch localElevator.Behaviour {
			case elevator.EB_Moving:
				if requests.RequestsShouldStop(localElevator) {
					fmt.Println("Elevator should stop")
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					fmt.Println("Door open")
					localElevator = requests.RequestsClearAtCurrentFloor(localElevator, hallRequestCleared)
					fmt.Println("Door open timer started")
					localElevator.Behaviour = elevator.EB_DoorOpen
					elevatorTx <- localElevator

					time.Sleep(3 * time.Second)
					setAllLights(localElevator)
					localElevator = onDoorTimeout(localElevator, hallRequestCleared)
					elevio.SetMotorDirection(localElevator.Dirn)
					elevatorTx <- localElevator
				}

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

		case obstruction := <-chObstructionSwitch:
			if obstruction {
				elevio.SetMotorDirection(elevio.MD_Stop)
				localElevator.Behaviour = elevator.EB_DoorOpen
				elevatorTx <- localElevator
			} else {
				elevio.SetMotorDirection(localElevator.Dirn)
				localElevator.Behaviour = elevator.EB_Moving
				elevatorTx <- localElevator
			}
		case <-chStopButton:
			//set behavior
			for f := 0; f < 4; f++ {
				for btn := elevio.ButtonType(0); btn < 3; btn++ {
					elevio.SetButtonLamp(btn, f, false)
				}
			}
		}

	}
}
