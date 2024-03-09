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
	if activeElevators[id-1] == true && id != 0 {
		return false
	} else {
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
	if id == 0 {
		masterState = true
	} else {
		masterState = false
	}
	fmt.Println("Starting elevator", id, "masterState:", masterState)

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

	go bcast.Transmitter(2000, elevatorTx)
	go bcast.Receiver(2001, elevatorRx)

	go bcast.Transmitter(3000, hallRequestsTx)
	go bcast.Receiver(3001, hallRequestsRx)

	go bcast.Transmitter(4000, watchdogTx)
	go bcast.Receiver(4001, watchdogRx)

	go watchdog.WatchdogCheckAlive(watchdogRx, chActiveWatchdogs, 100)
	go watchdog.WatchdogSendAlive(id, watchdogTx)

	go runElevator.RunLocalElevator(chActiveElevators, elevatorTx, hallRequestsRx, id)

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
		}
	}
}
