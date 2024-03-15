package main

import (
	"elevatorlib/elevator"
	"elevatorlib/elevator/runElevator"
	"elevatorlib/elevio"
	"elevatorlib/network/bcast"
	"elevatorlib/network/peers"
	"elevatorlib/requestAsigner"
	"flag"
	"fmt"
	"sort"
)

/*
	func checkMaster(chMasterState chan bool, id int, chActiveWatchdogs chan [3]bool) {
		fmt.Println("activeElevators changed checking master")
		oldMasterState := false
		newMasterState := oldMasterState
		for {
			select {
			case temp := <-chMasterState:
				oldMasterState = temp
			case activeElevators := <-chActiveWatchdogs:
				if id != 0 {
					if activeElevators[id-1] || activeElevators[0] {
						fmt.Println("Elevator:", id, "is slave")
						newMasterState = false
						fmt.Println("activeElevators", activeElevators)
					} else {
						fmt.Println("Elevator:", id, "is master")
						fmt.Println("activeElevators", activeElevators)
						newMasterState = true
					}
				} else {
					fmt.Println("Elevator:", id, "is master")
					fmt.Println("activeElevators", activeElevators)
					newMasterState = true
				}
				if oldMasterState != newMasterState {
					chMasterState <- newMasterState
					oldMasterState = newMasterState
				}
			default:
				fmt.Print(oldMasterState)
				chMasterState <- oldMasterState
			}

		}
	}
*/

func chooseMaster(peers []string) string {
	sort.Strings(peers)
	return peers[0]
}

type hallRequests map[string][][2]int

func main() {
	fmt.Println("Starting main")
	id := "0"
	port := 15657
	flag.StringVar(&id, "id", "0", "id of this elevator")
	flag.IntVar(&port, "port", 15657, "port of this elevator")
	flag.Parse()

	//check master state based on flag input
	//chanels
	chMasterState := make(chan bool)

	chElevatorTx := make(chan elevator.Elevator)
	chElevatorRx := make(chan elevator.Elevator)

	elevatorStatuses := make([]elevator.Elevator, 3)
	chElevatorStatuses := make(chan []elevator.Elevator)

	//Used for sending hall request too elevators from request assigner
	chAssignedHallRequestsTx := make(chan requestAsigner.HallRequests)
	chAssignedHallRequestsRx := make(chan requestAsigner.HallRequests)

	//used for sending new hall requests from elevators to request assigner
	chNewHallRequestTx := make(chan elevio.ButtonEvent)
	chNewHallRequestRx := make(chan elevio.ButtonEvent)
	//used to send information about cleared hall requests from elevators to request assigner
	chHallRequestClearedTx := make(chan elevio.ButtonEvent)
	chHallRequestClearedRx := make(chan elevio.ButtonEvent)

	chPeerEnable := make(chan bool)
	chPeerRxTx := make(chan peers.PeerUpdate)

	//used for sending elevator alive signal to watchdog
	//chWatchdogTx := make(chan int)
	//chWatchdogRx := make(chan int)
	//activeWatchdogs := [3]bool{false, false, false}
	//chActiveWatchdogs := make(chan [3]bool)

	//used for updating the active watchdogs array (checking which elevators are still alive)
	fmt.Println("Starting broadcast of, elevator, hallRequest and watchdog")
	//transmitter and receiver for elevator states
	go bcast.Transmitter(2000, chElevatorTx)
	go bcast.Receiver(2000, chElevatorRx)
	//transmitter and receiver for assigned hall requests
	go bcast.Transmitter(3001, chAssignedHallRequestsTx)
	go bcast.Receiver(3001, chAssignedHallRequestsRx)
	//transmitter and receiver for local hall requests
	go bcast.Transmitter(3002, chNewHallRequestTx)
	go bcast.Receiver(3002, chNewHallRequestRx)
	//transmitter and receiver for cleared hall requests
	go bcast.Transmitter(3003, chHallRequestClearedTx)
	go bcast.Receiver(3003, chHallRequestClearedRx)
	//transmitter and receiver for watchdog
	//go bcast.Transmitter(4001, chWatchdogTx)
	//go bcast.Receiver(4001, chWatchdogRx)
	go peers.Transmitter(4001, id, chPeerEnable)
	go peers.Receiver(4001, chPeerRxTx)
	//functions for checking the watchdog and sending alive signal
	//go watchdog.WatchdogSendAlive(id, chWatchdogTx)
	//go watchdog.WatchdogCheckAlive(chWatchdogRx, chActiveWatchdogs)

	//functions for running the local elevator
	go runElevator.RunLocalElevator(chElevatorTx, chNewHallRequestTx, chAssignedHallRequestsRx, chHallRequestClearedTx, id, port)

	//function for assigning hall request to slave elevators
	go requestAsigner.RequestAsigner(chNewHallRequestRx, chElevatorStatuses, chMasterState, chHallRequestClearedRx, chAssignedHallRequestsTx) //jobbe med den her

	fmt.Println("Starting peer")
	chPeerEnable <- true

	fmt.Println("Starting main loop")
	for {
		select {
		case elevator := <-chElevatorRx:
			elevatorStatuses[elevator.Id] = elevator
			chElevatorStatuses <- elevatorStatuses

		//case  <-chActiveWatchdogs:
		case pUpdate := <-chPeerRxTx:
			fmt.Println("Peers updated", pUpdate)
			if len(pUpdate.Peers) > 0 {
				masterID := chooseMaster(pUpdate.Peers)
				if id == masterID {
					fmt.Println("i am master")
					chMasterState <- true
				} else {
					fmt.Println("i am slave")
					chMasterState <- false
				}
			}
		}

	}
}
