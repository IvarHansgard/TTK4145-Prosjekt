package main

import (
	"elevatorlib/elevator"
	"elevatorlib/elevator/runElevator"
	"elevatorlib/elevio"
	"elevatorlib/network/bcast"
	"elevatorlib/requestAsigner"
	"elevatorlib/watchdog"
	"flag"
	"fmt"
)

func checkMaster(chMasterState chan bool, id int, activeElevators [3]bool) {
	fmt.Println("activeElevators changed checking master")
	if id != 0 {
		if activeElevators[id-1] || activeElevators[0] {
			fmt.Println("Elevator:", id, "is slave")
			chMasterState <- false
		} else {
			fmt.Println("Elevator:", id, "is master")
			chMasterState <- true
		}
	} else {
		fmt.Println("Elevator:", id, "is master")
		chMasterState <- true
	}
	return
}

type hallRequests map[string][][2]int

func main() {
	fmt.Println("Starting main")
	id := 0
	port := 15657
	flag.IntVar(&id, "id", 0, "id of this elevator")
	flag.IntVar(&port, "port", 15657, "port of this elevator")

	flag.Parse()
	//check master state based on flag input
	//chanels
	chMasterState := make(chan bool)

	elevatorTx := make(chan elevator.Elevator)
	elevatorRx := make(chan elevator.Elevator)
	elevatorStatuses := make([]elevator.Elevator, 3)
	chElevatorStatuses := make(chan []elevator.Elevator)

	assignedHallRequestsTx := make(chan requestAsigner.HallRequests)
	assignedHallRequestsRx := make(chan requestAsigner.HallRequests)
	chHallRequestClearedTx := make(chan elevio.ButtonEvent)
	chHallRequestClearedRx := make(chan elevio.ButtonEvent)
	chNewHallRequestTx := make(chan elevio.ButtonEvent)
	chNewHallRequestRx := make(chan elevio.ButtonEvent)

	chWatchdogTx := make(chan int)
	chWatchdogRx := make(chan int)
	chActiveWatchdogs := make(chan [3]bool)
	if id == 0 {
		chActiveWatchdogs <- [3]bool{true, false, false}
	} else {
		chActiveWatchdogs <- [3]bool{false, false, false}
	}

	fmt.Println("Starting broadcast of, elevator, hallRequest and watchdog")
	//transmitter and receiver for elevator states
	go bcast.Transmitter(2000, elevatorTx)
	go bcast.Receiver(2000, elevatorRx)

	//transmitter and receiver for assigned hall requests
	go bcast.Transmitter(3001, assignedHallRequestsTx)
	go bcast.Receiver(3001, assignedHallRequestsRx)
	//transmitter and receiver for local hall requests
	go bcast.Transmitter(3002, chNewHallRequestTx)
	go bcast.Receiver(3002, chNewHallRequestRx)
	//transmitter and receiver for cleared hall requests
	go bcast.Transmitter(3003, chHallRequestClearedTx)
	go bcast.Receiver(3003, chHallRequestClearedRx)

	//transmitter and receiver for watchdog
	go bcast.Transmitter(4001, chWatchdogTx)
	go bcast.Receiver(4001, chWatchdogRx)

	//functions for checking the watchdog and sending alive signal
	go watchdog.WatchdogSendAlive(id, chWatchdogTx)
	go watchdog.WatchdogCheckAlive(chWatchdogRx, chActiveWatchdogs, 10)

	//functions for running the local elevator
	go runElevator.RunLocalElevator(elevatorTx, chNewHallRequestTx, assignedHallRequestsRx, chHallRequestClearedTx, id, port)

	//function for assigning hall request to slave elevators
	go requestAsigner.RequestAsigner(chNewHallRequestRx, chElevatorStatuses, chMasterState, chHallRequestClearedRx, assignedHallRequestsTx) //jobbe med den her

	fmt.Println("Starting main loop")
	for {
		select {
		case elevator := <-elevatorRx:
			elevatorStatuses[elevator.Id] = elevator
			chElevatorStatuses <- elevatorStatuses

		case activeWatchdogs := <-chActiveWatchdogs:
			for i := 0; i < 3; i++ {
				if !activeWatchdogs[i] {
					elevatorStatuses[i].Behaviour = elevator.EB_Disconnected
				}
			}
			go checkMaster(chMasterState, id, activeWatchdogs)
		}
	}
}
