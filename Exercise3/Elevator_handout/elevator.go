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

	for {
		state: master
		nettverk()
		-The unassigned request
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
