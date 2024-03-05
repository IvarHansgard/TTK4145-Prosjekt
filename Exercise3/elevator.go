package main

import (
	"Driver-go/elevio"
	"fmt"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)
	connect_network()
	
	var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)


	/*func readUDP(port string, buffer chan []byte, n chan int) {

		recieveBuffer := make([]byte, 1024)
	
		udpAddress, err := net.ResolveUDPAddr("udp", port)
	
		udpConnection, err := net.ListenUDP("udp", udpAddress)
	
		defer udpConnection.Close() 
	
		numBytes, ipAddr, err := udpConnection.ReadFromUDP(recieveBuffer)
	
		if err != nil {
			log.Fatal(err, numBytes)
		}
	
		Println("you recieved: ", numBytes, "bytes from IP: ", ipAddr)
	
		buffer <- recieveBuffer
		n <- numBytes
	}*/
	
	func recieve_Elevator(port string) (Elevator, int, error){ //tar inn direction, floor og behaviour(door, moving, idle)
		recieveBuffer:= make([]byte, 1024)

		udpAddress, err:= net.ResolveUDPAddr("udp", port)
		udpConnection, err:=net.ListenUDP("udp", udpAddress)

		if err!=nil{
			return Elevator{}, 0, err
		}
		
		defer udpConnection.Close()
		numBytes, _, err:= udpConnection.ReadFromUDP(recieveBuffer)
		if err!= nil{
			return Elevator{}, 0, err
		}

		var myData Elevator
		buf:=bytes.NewReader(recieveBuffer[:numBytes])
		if err:= binary.Read(buf, binary.BigEndian, &myData); err!= nil{
			return Elevator{}, 0, err
		}
		return myData, numBytes, nil //returnerer hele structen Elevator 
		//for å vite hvilken heis, kan man tildele en port til per heis og derfor aktivt lete etter info fra en heis ved å sende inn tilsvarende port 
	}

func unassigned_requests(port string) (int){//hente ut tabell med alle requests. Alle fullførte requests fjernes fra tabellen og legges til i privat kø
	elevator:=recieve_Elevator(port)
	for i:=
}


	for {
		state: master
		nettverk()
		-The unassigned request//the hall requests
		-The whereabouts of the elevators (floor, direction, state/behavior (ie. moving, doorOpen, idle))
		-The current set of existing requests (cab requests and hall requests)
		-The availability or failure modes of the elevators (sara)
		elevator_algo() ivar
		send_queues()
		goto_floor() ulrikke

		state: slave
		get_cab_calls()
		send_hall_calls()
		get_queue()
		elevator_algo() ivar
		goto_floor()
		timeout_goto_mastermode()
		}
	}
}

