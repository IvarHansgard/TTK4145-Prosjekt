package main

//legge til alle imports
//fikse fsm så den er lik vanlig fsm

import (
	"ElevatorLib/driver-go/elevio"
	"ElevatorLib/network/bcast"
)

func main() {
	go elevio.Init("localhost:15657", 4)

	activeElevators := make(chan [3]bool)
	watchdogTx := make(chan int)      //elevator alive signal sent used when in slave and master mode
	watchdogRx := make(chan int)      //elevator alive signal recieved used when in slave and master mode
	elevatorTx := make(chan Elevator) //elevator objects sent used when in slave mode
	elevatorRx := make(chan Elevator) //elevator objects recieved used when in slave mode

	/*
		elevator1Rx := make(chan Elevator)  //elevator object 1 recieved used when in master mode
		elevator2Rx := make(chan Elevator)  //elevator object 2 recieved used when in master mode
		elevator1_aliveRx := make(chan int) //elevator alive signal 1 recieved used when in master mode
		elevator2_aliveRx := make(chan int) //elevator alive signal 2 recieved used when in master mode
	*/

	id := 0 //0, 1 or 2

	go bcast.Receiver(1000, watchdogRx) //
	go bcast.Receiver(2000, elevatorRx) //

	go bcast.Transmitter(1000, watchdogTx) //elevator 0: 3000 elevator 1: 3001 elevator 2: 3002
	go bcast.Transmitter(2000, elevatorTx) //elevator 0: 3000 elevator 1: 3001 elevator 2: 3002

	//run watchdog
	go watchdog.watchdog_sendAlive(id, watchdogTx)
	go watchdog.watchdog_checkAlive(watchdogRx, activeElevators, 10)

	/*
		//master
		go bcast.Receiver(2000, elevator0_aliveRx) //elevator 0: 2000 elevator 1: 2001 elevator 2: 2002
		go bcast.Receiver(2001, elevator1_aliveRx) //elevator 0: 2000 elevator 1: 2001 elevator 2: 2002
		go bcast.Receiver(2002, elevator2_aliveRx) //elevator 0: 2000 elevator 1: 2001 elevator 2: 2002

		go bcast.Receiver(3001, elevator1Rx) //elevator 0: 3000 elevator 1: 3001 elevator 2: 3002
		go bcast.Receiver(3002, elevator2Rx) //elevator 0: 3000 elevator 1: 3001 elevator 2: 3002

		go bcast.Transmitter(4000+id+1, elevator1Tx) //elevator 0: 3000 elevator 1: 3001 elevator 2: 3002
		go bcast.Transmitter(4000+id+2, elevator2Tx) //elevator 0: 3000 elevator 1: 3001 elevator 2: 3002
		//

		//slave
		go bcast.Receiver(1000+id, elevatorRx)    //elevator 0: 1000 elevator 1: 1001 elevator 2: 1002
		go bcast.Transmitter(2000+id, aliveTx)    //elevator 0: 2000 elevator 1: 2001 elevator 2: 2002
		go bcast.Transmitter(3000+id, elevatorTx) //elevator 0: 3000 elevator 1: 3001 elevator 2: 3002
		//
	*/

	id1 := -1
	id2 := -1
	timeOut := make(chan int)
	go watchdog(elevator1_aliveRx, timeOut, 10) //timeout for elevator1
	go watchdog(elevator2_aliveRx, timeOut, 10) //timeout for elevator2
	go watchdog_sendAlive(id, aliveTx)

	for {
		select {
		case a := <-elevator1_aliveRx:
			id1 = a.id

		case a := <-elevator2_aliveRx:
			id2 = a.id

		case <-timeoutSignal:
			break
		}
		if id1 && id2 != -1 {
			break
		}
	}

	master := false

	if id < id1 && id < id2 {
		master = true
	}

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	//queue moduletimeout
	hallRequests := make(chan [4][2]int)

	fsm_init(elevator1Rx, elevator2Rx)

	if master == true {
		for {
			select {
			case a := <-elevator1Rx:
				for x; x < 4; x++ {
					for y; y < 2; y++ {
						if a.requests[x][y] == 1 {
							hallRequests[x][y] = a.requests[x][y]
						}
					}
				}

			case a := <-elevator2Rx:
				for x; x < 4; x++ {
					for y; y < 2; y++ {
						if a.requests[x][y] == 1 {
							hallRequests[x][y] = a.requests[x][y]
						}
					}
				}

			case a := <-drv_buttons:
				fsm_onRequestButtonPress_master(a.Floor, a.Button, hallRequests)

			case a := <-drv_floors:
				onFloorArrival(a)

			case a := <-drv_stop:
				elevio.SetMotorDirection(elevio.MD_Stop)
			}
			fsm_run_algo()
			fsm_sendData_master(2, elevator1Tx)
			fsm_seneData_master(3, elevator2Tx)
		}
	} else {
		for {
			select {
			case a := <-drv_buttons:
				fsm_onRequestButtonPress_slave(a.Floor, a.Button)
			case a := <-drv_floors:
				fsm_onFloorArival()
			case a := <-drv_stop:
				elevio.SetMotorDirection(elevio.MD_Stop)
			case <-timeOut:
				fsm_change_to_master()
			case a := <-elevatorRx:
				fsm_recieveData(a)
			}
			fsm_sendData_slave(elevatorTx)
		}
	}
}
