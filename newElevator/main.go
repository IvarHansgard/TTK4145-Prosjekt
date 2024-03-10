package main

import (
	"elevatorlib/elevator"
	"elevatorlib/elevator/runElevator"
	"elevatorlib/network/bcast"
	"elevatorlib/requestAsigner"
	"elevatorlib/watchdog"
	"flag"
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
	id := 0
	port := 15657
	flag.IntVar(&id, "id", 0, "id of this elevator")
	flag.IntVar(&port, "port", 15657, "port of this elevator")

	flag.Parse()

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

	assignedHallRequestsTx := make(chan requestAsigner.HallRequests)
	assignedHallRequestsRx := make(chan requestAsigner.HallRequests)
	localHallRequestsTx := make(chan [][2]bool)
	localHallRequestsRx := make(chan [][2]bool)
	chHallRequestClearedTx := make(chan [2]int)
	chHallRequestClearedRx := make(chan [2]int)

	watchdogTx := make(chan int)
	watchdogRx := make(chan int)
	chActiveWatchdogs := make(chan []bool)

	fmt.Println("Starting broadcast of, elevator, hallRequest and watchdog")
	//transmitter and receiver for elevator states
	go bcast.Transmitter(2000, elevatorTx)
	go bcast.Receiver(2000, elevatorRx)

	//transmitter and receiver for assigned hall requests
	go bcast.Transmitter(3001, assignedHallRequestsTx)
	go bcast.Receiver(3001, assignedHallRequestsRx)
	//transmitter and receiver for local hall requests
	go bcast.Transmitter(3002, localHallRequestsTx)
	go bcast.Receiver(3002, localHallRequestsRx)
	//transmitter and receiver for cleared hall requests
	go bcast.Transmitter(3003, chHallRequestClearedTx)
	go bcast.Receiver(3003, chHallRequestClearedRx)

	//transmitter and receiver for watchdog
	go bcast.Transmitter(4001, watchdogTx)
	go bcast.Receiver(4001, watchdogRx)

	//functions for checking the watchdog and sending alive signal

	go watchdog.WatchdogCheckAlive(watchdogRx, chActiveWatchdogs, 10)
	go watchdog.WatchdogSendAlive(id, watchdogTx)

	//functions for running the local elevator
	go runElevator.RunLocalElevator(chActiveElevators, elevatorTx, localHallRequestsTx, assignedHallRequestsRx, chHallRequestClearedTx, id, port)

	//function for assigning hall request to slave elevators

	go requestAsigner.RequestAsigner(chActiveElevators, masterState, localHallRequestsRx, assignedHallRequestsTx, chHallRequestClearedRx) //jobbe med den her

	fmt.Println("Starting main loop")
	for {
		select {
		case elevator := <-elevatorRx:
			activeElevators[elevator.Id] = elevator
			chActiveElevators <- activeElevators

		case activeWatchdogs := <-chActiveWatchdogs:
			for i := 0; i < 3; i++ {
				if !activeWatchdogs[i] {
					activeElevators[i].Behaviour = elevator.EB_Disconnected
				}
			}
			masterState = checkMaster(id, activeWatchdogs)
			chMasterState <- masterState
			//assign lost elevators orders to other elevators
		}
	}
}
