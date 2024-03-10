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

func onDoorTimeout(e elevator.Elevator) elevator.Elevator {
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
			e = requests.RequestsClearAtCurrentFloor(e)
			fmt.Println(e)
			onDoorTimeout(e)
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

	e.Dirn = requests.RequestsChooseDirection(e).Dirn
	e.Behaviour = requests.RequestsChooseDirection(e).Behaviour
	elevio.SetMotorDirection(e.Dirn)

	fmt.Println("MotorDirection: ", e.Dirn)

	return e
}

func localElevatorInit(id int) elevator.Elevator {
	fmt.Println("Initializing elevator ", id)

	elevio.Init("localhost:15657", 4)
	elevator := elevator.Elevator_init(id)
	setAllLights(elevator)
	elevio.SetDoorOpenLamp(false)

	return elevator
}

func RunLocalElevator(chActiveElevators chan []elevator.Elevator, elevatorTx chan elevator.Elevator, hallRequestsRx chan requestAsigner.HallRequests, id int) {
	fmt.Println("Starting localElevator")
	localElevator := localElevatorInit(id)

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

			fmt.Println("Button press detected")
			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				fmt.Println("Door is open")
				if requests.RequestsShouldClearImmediately(localElevator, button.Floor, button.Button) {
					time.Sleep(3 * time.Second)
					localElevator = onDoorTimeout(localElevator)
				} else {
					localElevator.Requests[button.Floor][button.Button] = true
				}

			case elevator.EB_Moving:
				fmt.Println("Elevator is moving.")
				localElevator.Requests[button.Floor][button.Button] = true

			case elevator.EB_Idle:
				fmt.Println("Elevator is idle.")
				localElevator.Requests[button.Floor][button.Button] = true
				pair := requests.RequestsChooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour
				switch pair.Behaviour {
				case elevator.EB_DoorOpen:
					fmt.Println("Door is open")
					elevio.SetDoorOpenLamp(true)
					time.Sleep(3 * time.Second)
					localElevator = onDoorTimeout(localElevator)
					localElevator = requests.RequestsClearAtCurrentFloor(localElevator)
				case elevator.EB_Moving:
					fmt.Println("Elevator is moving.")
					elevio.SetMotorDirection(localElevator.Dirn)
				case elevator.EB_Idle:
					break
				}
			}
			setAllLights(localElevator)
			elevatorTx <- localElevator //check hall request ikke oppdater nå hall

		case floor := <-chFloor:
			localElevator.Floor = floor
			elevio.SetFloorIndicator(localElevator.Floor)

			switch localElevator.Behaviour {
			case elevator.EB_Moving:
				if requests.RequestsShouldStop(localElevator) {
					fmt.Println("Elevator should stop")
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					fmt.Println("Door open")
					localElevator = requests.RequestsClearAtCurrentFloor(localElevator)

					fmt.Println("Door open timer started")
					localElevator.Behaviour = elevator.EB_DoorOpen
					elevatorTx <- localElevator

					time.Sleep(3 * time.Second)
					setAllLights(localElevator)
					localElevator = onDoorTimeout(localElevator)
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
							elevatorTx <- localElevator
							break
						}
					}
				}

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
