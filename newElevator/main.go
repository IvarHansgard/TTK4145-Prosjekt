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
)

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
	chPeerRx := make(chan peers.PeerUpdate)

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
	go runElevator.RunLocalElevator(chElevatorTx, chNewHallRequestTx, chAssignedHallRequestsRx, chHallRequestClearedTx, id, port)

	//function for assigning hall request to slave elevators
	go requestAsigner.RequestAsigner(chNewHallRequestRx, chElevatorStatuses, chRequestAssignerMasterState, chHallRequestClearedRx, chAssignedHallRequestsTx) //jobbe med den her

	fmt.Println("Starting main loop")
	for {
		select {
		case elevator := <-chElevatorRx:
			elevatorStatuses[elevator.Id] = elevator
			chElevatorStatuses <- elevatorStatuses

		case pUpdate := <-chPeerRx:
			fmt.Printf("Peer update:\n")
			go checkMaster(chMasterState, masterState, id, pUpdate)

		case state := <-chMasterState:
			masterState = state
			chRequestAssignerMasterState <- masterState
		}
	}
}
