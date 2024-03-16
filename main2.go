package main

import (
	"elevatorlib/network/peers"
	"flag"
	"fmt"
)

func main() {
	id := "0"
	flag.StringVar(&id, "id", "0", "id of this elevator")
	flag.Parse()

	masterMode := true

	peerChan := make(chan peers.PeerUpdate)
	peerEnable := make(chan bool)
	fmt.Println("Starting reciever")
	go peers.Receiver(1000, peerChan)
	fmt.Println("Starting transmitter")
	go peers.Transmitter(1000, id, peerEnable)

	for {
		select {
		case peerUpdate := <-peerChan:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peerUpdate.Peers)
			fmt.Printf("  New:      %q\n", peerUpdate.New)
			fmt.Printf("  Lost:     %q\n", peerUpdate.Lost)
			if len(peerUpdate.Peers) == 1 && peerUpdate.Peers[0] != id {
				fmt.Println("only one peer that is not me ignoring")
			} else if len(peerUpdate.Peers) == 1 && peerUpdate.Peers[0] == id {
				fmt.Println("I am the only peer")
				fmt.Println("")
				if len(peerUpdate.Lost) > 0 && peerUpdate.Peers[0] == id {
					fmt.Println("Im master")
					masterMode = true
				}
			} else if peerUpdate.Peers[0] == id {
				if !masterMode {
					masterMode = true
					fmt.Println("I am the new master")
				} else {
					fmt.Println("I am still master")
				}
			} else {
				if masterMode {
					masterMode = false
					fmt.Println("I am no longer master")
				} else {
					fmt.Println("I am still slave")
				}
			}

		case enable := <-peerEnable:
			fmt.Println("enable", enable)
		}
	}
}
