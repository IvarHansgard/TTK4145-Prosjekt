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
	"strconv"
)

/*
	func checkMaster(chMasterState chan bool, masterState bool, id string, pUpdate peers.PeerUpdate) {
		if len(pUpdate.Peers) == 1 && pUpdate.New == id {
			fmt.Println("Start Up")
			fmt.Println("My master state is", masterState)
		} else if len(pUpdate.Lost) > 0 {
			fmt.Println("Lost peer", pUpdate.Lost)
			fmt.Println("checking master")
			if pUpdate.Peers[0] == id {
				if !masterState {
					fmt.Println("Changing to master")
					chMasterState <- true
				}
				fmt.Println("I am master")
			} else {
				if masterState {
					fmt.Println("Changing to slave")
					chMasterState <- false
				}
				fmt.Println("I am slave")
			}
		} else if pUpdate.New != "" && pUpdate.New != id {
			fmt.Println("New peer", pUpdate.New)
			fmt.Println("checking master")
			if pUpdate.Peers[0] == id {
				if !masterState {
					fmt.Println("Changing to master")
					chMasterState <- true
				}
				fmt.Println("I am master")
			} else {
				if masterState {
					fmt.Println("Changing to slave")
					chMasterState <- false
				}
				fmt.Println("I am slave")
			}
		} else if len(pUpdate.Peers) > 0 {
			if pUpdate.Peers[0] == id {
				if !masterState {
					fmt.Println("Changing to master")
					chMasterState <- true
				}
				fmt.Println("I am master")
			} else {
				if masterState {
					fmt.Println("Changing to slave")
					chMasterState <- false
				}
				fmt.Println("I am slave")
			}
		}
	}
*/

/*
	func checkMaster(chMasterState chan bool, masterState bool, id string, peerUpdate peers.PeerUpdate) {
		if len(peerUpdate.Peers) == 1 && peerUpdate.New == id {
			fmt.Println("Start Up")
			fmt.Println("My master state is", masterState)
		} else if len(peerUpdate.Lost) > 0 {
			fmt.Println("Lost peer", peerUpdate.Lost)
			fmt.Println("checking master")
			if peerUpdate.Peers[0] == id {
				if !masterState {
					fmt.Println("Changing to master")
					chMasterState <- true
				}
				fmt.Println("I am master")
			} else {
				if masterState {
					fmt.Println("Changing to slave")
					chMasterState <- false
				}
				fmt.Println("I am slave")
			}
		} else if peerUpdate.New != "" && peerUpdate.New != id {
			fmt.Println("New peer", peerUpdate.New)
			fmt.Println("checking master")
			if peerUpdate.Peers[0] == id {
				if !masterState {
					fmt.Println("Changing to master")
					chMasterState <- true
				}
				fmt.Println("I am master")
			} else {
				if masterState {
					fmt.Println("Changing to slave")
					chMasterState <- false
				}
				fmt.Println("I am slave")
			}
		} else if len(peerUpdate.Peers) > 0 {
			if peerUpdate.Peers[0] == id {
				if !masterState {
					fmt.Println("Changing to master")
					chMasterState <- true
				}
				fmt.Println("I am master")
			} else {
				if masterState {
					fmt.Println("Changing to slave")
					chMasterState <- false
				}
				fmt.Println("I am slave")
			}
		}
	}
*/
func assignDisconnected(id string, peerUpdate peers.PeerUpdate, elevatorStates []elevator.Elevator, chElevatorStates chan []elevator.Elevator, chElevatorLost chan bool) {
	if len(peerUpdate.Lost) > 0 && peerUpdate.Lost[0] != id {
		for i := 0; i < len(peerUpdate.Lost); i++ {
			elevId, err := strconv.Atoi(peerUpdate.Lost[i])
			if err != nil {
				fmt.Println("Error:", err)
			}
			elevatorStates[elevId].Behaviour = elevator.EB_Disconnected
			elevatorStates[elevId].Dirn = elevio.MD_Stop
			fmt.Println("elevator", elevId, "disconnected")
			chElevatorStates <- elevatorStates
			chElevatorLost <- true
		}
	}
	return
}

func checkMaster(chMasterState chan bool, masterState bool, id string, peerUpdate peers.PeerUpdate) {
	ignore := false
	for i := 0; i < len(peerUpdate.Lost); i++ {
		if peerUpdate.Lost[i] == id {
			ignore = true
		}
	}
	if peerUpdate.New == id {
		ignore = true
	}
	if !ignore {
		fmt.Printf("Peer update:\n")
		fmt.Printf("  Peers:    %q\n", peerUpdate.Peers)
		fmt.Printf("  New:      %q\n", peerUpdate.New)
		fmt.Printf("  Lost:     %q\n", peerUpdate.Lost)
		if peerUpdate.Peers[0] == id {
			if !masterState {
				fmt.Println("changed master state to true")
				chMasterState <- true
			}
		} else {
			if masterState {
				fmt.Println("changed master state to false")
				chMasterState <- false
			}
		}
	}
	return
}

type hallRequests map[string][][2]int

func main() {
	fmt.Println("Starting main")
	id := "0"
	port := 15657
	flag.StringVar(&id, "id", "0", "id of this elevator")
	flag.IntVar(&port, "port", 15657, "port of this elevator")
	flag.Parse()

	//chanels
	masterState := true
	chMasterState := make(chan bool)
	chRequestAssignerMasterState := make(chan bool)

	chElevatorTx := make(chan elevator.Elevator)
	chElevatorRx := make(chan elevator.Elevator)

	elevatorStates := make([]elevator.Elevator, 3)
	chElevatorStates := make(chan []elevator.Elevator)

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
	chPeerRx := make(chan peers.PeerUpdate)

	chStopButtonPressed := make(chan bool)
	chElevatorLost := make(chan bool)

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
	//transmitter and receiver for peer
	go peers.Transmitter(4001, id, chPeerEnable)
	go peers.Receiver(4001, chPeerRx)

	//functions for running the local elevator
	go runElevator.RunLocalElevator(chElevatorTx, chNewHallRequestTx, chAssignedHallRequestsRx, chHallRequestClearedTx, id, port, chStopButtonPressed)

	//function for assigning hall request to slave elevators
	go requestAsigner.RequestAsigner(chNewHallRequestRx, chElevatorStates, chRequestAssignerMasterState, chHallRequestClearedRx, chAssignedHallRequestsTx,chStopButtonPressed) //jobbe med den her

	fmt.Println("Starting main loop")
	for {
		select {
		case elevator := <-chElevatorRx:
			elevatorStates[elevator.Id] = elevator
			chElevatorStates <- elevatorStates

		case peerUpdate := <-chPeerRx:
			fmt.Printf("Peer update:\n")
			go assignDisconnected(id, peerUpdate, elevatorStates, chElevatorStates, chElevatorLost)
			go checkMaster(chMasterState, masterState, id, peerUpdate)

		case state := <-chMasterState:
			masterState = state
			chRequestAssignerMasterState <- masterState
		}
	}
}
