package watchdog

import (
	"time"
)

func watchdog_checkAlive(elevatorSignal chan int, activeElevators chan bool timeout int) {
	initialTimeout := timeout
	timeoutElevator0 := timeout
	timeoutElevator1 := timeout
	timeoutElevator2 := timeout

	for {
		select {
		case elevator_id := <-elevatorSignal:
			if elevator_id == 0 {
				timeoutElevator0 = initialTimeout
			} else if elevator_id == 1 {
				timeoutElevator1 = initialTimeout
			} else if elevator_id == 2 {
				timeoutElevator2 = initialTimeout
			}
		}

		//decrease 
		timeoutElevator0--
		timeoutElevator1--
		timeoutElevator2--

		//check if any of the elevators have timed out and set the elevators active status to false
		if timeoutElevator0 == 0 {
			activeElevators[0] <- 0
		}
		if timeoutElevator1 == 0 {
			activeElevators[1] <- 0
		}
		if timeoutElevator2 == 0 {
			activeElevators[2] <- 0
		}

		time.Sleep(1 * time.Second)
	}
}

func watchdog_sendAlive(id int, tx chan int) {
	tx <- id
}
