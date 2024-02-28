package main

import "_/home/student/Documents/gruppe23/asdas/TTK4145-Prosjekt/Exercise3/driver-go/elevio"



func getID("id") int {
	return id
}

func sendID(id int){
	
}

func main() {
	go elevio.Init("localhost:15657", 4)
	
	id := 0	

	sendID(id)

	id1 := getID("ip")e1
	id2 := getID("ip")
	
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

	slave1 := make(chan Elevatoe1
	go recieve_elevator(port1, slave1)
	go recieve_elevator(port2, slave2)
	
	timer := make(chan int)
	go watchdog(timer)

	queue := make(chan int[4][3])

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
				watchdog <- 300

			case a := <- slave2:
				for x ; x < 4; x++{
					for y ; y < 2; y++{
						if(a.requests[x][y] == 1){
							hallRequests[x][y] = a.requests[x][y]
						}
					}		
				}
				watchdog <-300

			case a := <-drv_buttons:
				fsm_onRequestButtonPress(a.Floor, a.Button, hallRequests)

			case a := <-drv_floors:
				onFloorArrival(a)

			case a := <-drv_stop:
				elevio.SetMotorDirection(elevio.MD_Stop)
			}
			fsm_run_algo()
			fsm_check_requests()
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
			}
			fsm_sendData()
			fsm_recieveData()
			fsm_check_requests()
			fsm_check_watchdog()
			}
		}
}	
