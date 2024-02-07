package main

import (
	"Network-go/network/bcast"
	"Network-go/network/localip"
	"Network-go/network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//
//	will be received as zero-values.
type Elevator struct {
	Name  string
	Floor int
}

func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that aretype Elevator struct {
	Name  string
	Floor int
}
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(2001, id, peerTxEnable) //2001
	go peers.Receiver(2000, peerUpdateCh)        //2000

	// We make channels for sending and receiving our custom data types
	elevatorTx := make(chan Elevator)
	elevatorRx := make(chan Elevator)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(2003, elevatorTx) //2003
	go bcast.Receiver(2002, elevatorRx)    //2002

	// The example message. We just send one of these every second.
	go func() {
		E1 := Elevator{"Elevator-2", 0}
		for {
			E1.Floor++
			if E1.Floor >= 4 {
				E1.Floor = 0
			}
			elevatorTx <- E1
			time.Sleep(2 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-elevatorRx:
			fmt.Printf("Elevator name: %s \nFloor: %d\n", a.Name, a.Floor)
		}
	}
}
