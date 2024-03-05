package main

import(
	"fmt"
	"time"
	"ElevatorLib/elevator"
	"ElevatorLib/requests"
	"ElevatorLib/driver-go/elevio"
) 

//network package
func getID("id") int {
	return id
}

func sendID(id int){
	
}

func main() {
	go elevio.Init("localhost:15657", 4)
	
	id := 0	

	sendID(id)

	id1 := getID("ip")
	id2 := getID("ip")
	
	go getID("ip")
	go getID("ip")

	for(){
		select{
		case a := <- id1:
		case a := <- id2:
		}
		if id1 && id2 != -1{
			break
		}
	}
	master := false

	if id < id1 and id < id2{
		master = true
	}


	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)


	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	slave1 := make(chan Elevator)
	slave2 := make(chan Elevator)
	
	//network package
	go recieve_elevator(port1, slave1)
	go recieve_elevator(port2, slave2)
	

	//make watchdog function
	go recive_alive_signal(port1, slave1_alive)
	go recive_alive_signal(port2, slave2_alive)
	go watchdog(slave1_alive,slave2_alive, 300)

	//queue module
	hallRequests := make(chan int[4][2])

	if master == true{
	for {
		select {
			case a := <- slave1:
				for x ; x < 4; x++{
					for y ; y < 2; y++{
						if(a.requests[x][y] == 1){
							hallRequests[x][y] = a.requests[x][y]
						}
					}		
				}
				watchdog1 <- 300

			case a := <- slave2:
				for x ; x < 4; x++{
					for y ; y < 2; y++{
						if(a.requests[x][y] == 1){
							hallRequests[x][y] = a.requests[x][y]
						}
					}		
				}
				watchdog2 <- 300

			case a := <-drv_buttons:
				fsm_onRequestButtonPress_master(a.Floor, a.Button, hallRequests)

			case a := <-drv_floors:
				onFloorArrival(a)

			case a := <-drv_stop:
				elevio.SetMotorDirection(elevio.MD_Stop)
			}
			fsm_run_algo()
		}
	}else{
		for {
			select {
				case a := <-drv_buttons:
					fsm_onRequestButtonPress_slave(a.Floor, a.Button)
				case a := <-drv_floors:
					fsm_onFloorArival()
				case a := <-drv_stop:
					elevio.SetMotorDirection(elevio.MD_Stop)
				case a := <- watchdog:
				if a == 0{
					fsm_change_to_master()
				}
				}
			fsm_sendData()
			fsm_recieveData()
			requests_chooseDirection()

			}
		}
}	
