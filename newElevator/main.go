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

type hallRequests map[string][][2]int

func checkMasterState(peerUpdate peers.PeerUpdate, chMasterState chan bool, id string, chElevatorLost chan string) {
	masterState := true
	if len(peerUpdate.Peers) == 1 && peerUpdate.Peers[0] != id {
		//debug
		//fmt.Println("only one peer that is not me ignoring")
		//fmt.Println("master state", masterState)
	} else if len(peerUpdate.Peers) == 1 && peerUpdate.Peers[0] == id {
		fmt.Println("I am the only peer")
		if len(peerUpdate.Lost) > 0 && peerUpdate.Peers[0] == id {
			fmt.Println("Im master")
			masterState = true
			chMasterState <- true
		}
	} else if peerUpdate.Peers[0] == id {
		if !masterState {
			masterState = true
			chMasterState <- true
			fmt.Println("I am the new master")
		} else {
			fmt.Println("I am still master")
		}
	} else {
		if masterState {
			masterState = false
			chMasterState <- false
			fmt.Println("I am no longer master")
		} else {
			fmt.Println("I am still slave")
		}
	}
	if len(peerUpdate.Lost) > 0 {
		fmt.Println("peer lost", peerUpdate.Lost)
		if peerUpdate.Lost[0] != id {
			for i := 0; i < len(peerUpdate.Lost); i++ {
				chElevatorLost <- peerUpdate.Lost[i]
			}

		}
	}
	return
}

func main() {
	fmt.Println("Starting main")
	id := "0"
	port := 15657
	numElevators := 3
	flag.StringVar(&id, "id", "0", "id of this elevator")
	flag.IntVar(&port, "port", 15657, "port of this elevator")
	flag.IntVar(&numElevators, "elevators", 3, "numbers of elevators in the system")
	flag.Parse()

	//chanels
	chMasterState := make(chan bool)

	chElevatorTx := make(chan elevator.Elevator)
	chElevatorRx := make(chan elevator.Elevator)

	//Used for sending hall request too elevators from request assigner
	chAssignedHallRequestsTx := make(chan requestAsigner.HallRequests)
	chAssignedHallRequestsRx := make(chan requestAsigner.HallRequests)

	//used for sending new hall requests from elevators to request assigner
	chNewHallRequestTx := make(chan elevio.ButtonEvent)
	chNewHallRequestRx := make(chan elevio.ButtonEvent)
	chSendHallRequestsToMasterTx := make(chan [4][2]bool)
	chSendHallRequestsToMasterRx := make(chan [4][2]bool)
	chSendElevatorStatesToMasterTx := make(chan requestAsigner.ElevatorMap)
	chSendElevatorStatesToMasterRx := make(chan requestAsigner.ElevatorMap)
	//used to send information about cleared hall requests from elevators to request assigner
	chHallRequestClearedTx := make(chan elevio.ButtonEvent)
	chHallRequestClearedRx := make(chan elevio.ButtonEvent)
	chSetButtonLightTx := make(chan elevio.ButtonEvent)
	chSetButtonLightRx := make(chan elevio.ButtonEvent)

	chPeerEnable := make(chan bool)
	chPeerRx := make(chan peers.PeerUpdate)

	chStopButtonPressed := make(chan bool)
	chElevatorLost := make(chan string)

	//initing num of elevators

	//used for updating the active watchdogs array (checking which elevators are still alive)
	fmt.Println("Starting broadcast of, elevator, hallRequest and watchdog")
	//transmitter and receiver for elevator states
	go bcast.Transmitter(2000, chElevatorTx)
	go bcast.Receiver(2000, chElevatorRx)
	//transmitter and receiver for assigned hall requests
	go bcast.Transmitter(3000, chAssignedHallRequestsTx)
	go bcast.Receiver(3000, chAssignedHallRequestsRx)
	//transmitter and receiver for local hall requests
	go bcast.Transmitter(4000, chNewHallRequestTx)
	go bcast.Receiver(4000, chNewHallRequestRx)
	//
	go bcast.Transmitter(5000, chSendHallRequestsToMasterTx)
	go bcast.Receiver(5000, chSendHallRequestsToMasterRx)
	go bcast.Transmitter(6000, chSendElevatorStatesToMasterTx)
	go bcast.Receiver(6000, chSendElevatorStatesToMasterRx)
	//transmitter and receiver for cleared hall requests
	go bcast.Transmitter(7000, chHallRequestClearedTx)
	go bcast.Receiver(7000, chHallRequestClearedRx)
	//transmitter and receiver for setting button lights
	go bcast.Transmitter(8000, chSetButtonLightTx)
	go bcast.Receiver(8000, chSetButtonLightRx)
	//Peer transmitter and receiver
	fmt.Println("starting peers")
	go peers.Transmitter(9000, id, chPeerEnable)
	go peers.Receiver(9000, chPeerRx)

	//functions for running the local elevator
	go runElevator.RunLocalElevator(chElevatorTx, chNewHallRequestTx, chAssignedHallRequestsRx,
		chHallRequestClearedTx, id, port,
		chStopButtonPressed, chSetButtonLightRx, chSetButtonLightTx)

	//function for assigning hall request to slave elevators
	go requestAsigner.RequestAsigner(chNewHallRequestRx, chElevatorRx, chMasterState,
		chHallRequestClearedRx, chAssignedHallRequestsTx, chStopButtonPressed,
		chSendHallRequestsToMasterTx, chSendHallRequestsToMasterRx, chSendElevatorStatesToMasterTx,
		chSendElevatorStatesToMasterRx, chElevatorLost, numElevators)

	fmt.Println("Starting main loop")
	for {
		select {
		case peerUpdate := <-chPeerRx:
			//debug
			//fmt.Printf("Peer update:\n")
			//fmt.Println("Peers:", peerUpdate.Peers)
			//fmt.Println("New:", peerUpdate.New)
			go checkMasterState(peerUpdate, chMasterState, id, chElevatorLost)
		}
	}
}
