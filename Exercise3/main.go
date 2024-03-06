package main

//legge til alle imports
//fikse fsm så den er lik vanlig fsm

import (
	"ElevatorLib/driver-go/elevio"
	. "ElevatorLib/elevator"
	"ElevatorLib/fsm"
	"ElevatorLib/network/bcast"
	"ElevatorLib/watchdog"
)

func main() {
	go elevio.Init("localhost:15657", 4)

	master := false

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
	go watchdog.Watchdog_sendAlive(id, watchdogTx)
	go watchdog.Watchdog_checkAlive(watchdogRx, activeElevators, 10)

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
	*/[4][2]int
	temp := <-activeElevators
	if temp[id-1] == true && id != 0 {
		master = false
	} else {
		master = true
	}

	elevatorArray := make([]Elevator, 3)
	for i := 0; i < 3; i++ {
		elevatorArray[i] = Elevator_uninitialized()
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
	var temp_hallRequests [4][2]int

	fsm.Fsm_init()

	if master == true {
		for {
			select {
			case a := <-elevatorRx:
				for x := 0; x < 4; x++ {
					for y := 0; y < 2; y++ {
						if a.Requests[x][y] == 1 {
							temp_hallRequests[x][y] = a.Requests[x][y]
						}
					}
				}
				hallRequests <- temp_hallRequests

			case a := <-drv_buttons:
				fsm.Fsm_onRequestButtonPress_master(a.Floor, a.Button, hallRequests)

			case a := <-drv_floors:
				fsm.OnFloorArrival(a)

			case a := <-drv_stop:
				if a == true {
					elevio.SetMotorDirection(elevio.MD_Stop)
				} //fikk feilmelding fordi a ike blir brukt så jeg la til if a==true
			}

			fsm.Fsm_run_algo()

			fsm.Fsm_sendData_master(2, elevatorRx)
			fsm.Fsm_sendData_master(3, elevatorRx)
		}
	} else {
		for {
			select {
			case a := <-activeElevators:
				fsm.Fsm_check_master(a)

			case a := <-drv_buttons:
				fsm.Fsm_onRequestButtonPress_slave(a.Floor, a.Button)
			case a := <-drv_floors:
				fsm.OnFloorArrival(a)
			case a := <-drv_stop:
				if a == true {
					elevio.SetMotorDirection(elevio.MD_Stop)
				}
			case a := <-elevatorRx:
				fsm.Fsm_recieveData(a)
			}
			fsm.Fsm_sendData_slave(elevatorTx)
		}
	}
}
