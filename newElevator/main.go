package main

import (
	"elevatorlib/elevator"
	"elevatorlib/elevator/runElevator"
	"elevatorlib/network/bcast"
	"elevatorlib/requestAsigner"
	"elevatorlib/watchdog"
	"fmt"
)

func checkMaster(id int, activeElevators []bool) bool {
	fmt.Println("activeElevators changed checking master")
	if activeElevators[id-1] == true && id != 0 {
		fmt.Println("im not master")
		return false
	} else {
		fmt.Println("im master")
		return true
	}
}

type hallRequests map[string][][2]int

func main() {
	fmt.Println("Starting main")
	/*
		var id int
		flag.IntVar(&id, "id", 0, "id of this elevator")
		flag.Parse()
	*/
	id := 0
	var masterState bool

	//check master state based on flag input
	if id == 0 {
		masterState = true
	} else {
		masterState = false
	}

	//chanels
	chMasterState := make(chan bool)

	elevatorTx := make(chan elevator.Elevator)
	elevatorRx := make(chan elevator.Elevator)
	activeElevators := make([]elevator.Elevator, 3)
	chActiveElevators := make(chan []elevator.Elevator)

	hallRequestsTx := make(chan requestAsigner.HallRequests)
	hallRequestsRx := make(chan requestAsigner.HallRequests)

	watchdogTx := make(chan int)
	watchdogRx := make(chan int)
	chActiveWatchdogs := make(chan []bool)

	fmt.Println("Starting broadcast of, elevator, hallRequest and watchdog")
	//transmitter and receiver for elevator states
	go bcast.Transmitter(2000, elevatorTx)
	go bcast.Receiver(2001, elevatorRx)

	//transmitter and receiver for assigned hall requests
	go bcast.Transmitter(3000, hallRequestsTx)
	go bcast.Receiver(3001, hallRequestsRx)

	//transmitter and receiver for watchdog
	go bcast.Transmitter(4000, watchdogTx)
	go bcast.Receiver(4001, watchdogRx)

	//functions for checking the watchdog and sending alive signal

	go watchdog.WatchdogCheckAlive(watchdogRx, chActiveWatchdogs, 100)
	go watchdog.WatchdogSendAlive(id, watchdogTx)

	//functions for running the local elevator
	go runElevator.RunLocalElevator(chActiveElevators, elevatorTx, hallRequestsRx, id)

	//function for assigning hall request to slave elevators

	go requestAsigner.RequestAsigner(chActiveElevators, masterState, hallRequestsTx) //jobbe med den her

	fmt.Println("Starting main loop")
	for {
		select {
		case elevator := <-elevatorRx:
			activeElevators[elevator.Id] = elevator
			chActiveElevators <- activeElevators

		case activeWatchdogs := <-chActiveWatchdogs:
			masterState = checkMaster(id, activeWatchdogs)
			chMasterState <- masterState
			//assign lost elevators orders to other elevators
		}
	}
}
