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
	
	elevatorTx := make(chan Elevator)
	
	elevator1Rx := make(chan Elevator)
	elevator2Rx := make(chan Elevator)
	elevator1_aliveRx := make(chan int)
	elevator2_aliveRx := make(chan int)

	go bcast.Receiver(2001, elevator1_aliveRx)    //elevator 0: 2000 elevator 1: 2001 elevator 2: 2002
	go bcast.Receiver(2002, elevator2_aliveRx)    //elevator 0: 2000 elevator 1: 2001 elevator 2: 2002

	go bcast.Receiver(3001, elevator1Rx)    //elevator 0: 3000 elevator 1: 3001 elevator 2: 3002 
	go bcast.Receiver(3002, elevator2Rx)    //elevator 0: 3000 elevator 1: 3001 elevator 2: 3002
	
	id := 0	

	go bcast.Transmitter(1000+id, aliveTx) //1000
	go bcast.Transmitter(2000+id, elevatorTx) //2000

	id1 := -1
	id2 := -1

	for(){
		select{
		case a := <- elevator1_aliveRx:
			id1 = a.id
		case a := <- elevator2_aliveRx:
			id2 = a.id
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

	go watchdog(elevator1_aliveRx,elevator2_aliveRx, 300)

	//queue module
	hallRequests := make(chan int[4][2])

	if master == true{
	for {
		select {
			case a := <- elevator1Rx:
				for x ; x < 4; x++{
					for y ; y < 2; y++{
						if(a.requests[x][y] == 1){
							hallRequests[x][y] = a.requests[x][y]
						}
					}		
				}
				watchdog1 <- 300

			case a := <- elevator2Rx:
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
