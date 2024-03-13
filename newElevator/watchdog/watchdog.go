package watchdog

import (
	"fmt"
	"time"
)

func WatchdogCheckAlive(elevatorSignal chan int, activeWatchdogs chan [3]bool, timeout int) {
	fmt.Println("Starting watchdog check alive")
	prevTemp := [3]bool{false, false, false}
	temp := [3]bool{false, false, false}

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
 
		//decrease timeout

		timeoutElevator0--

		timeoutElevator1--

		timeoutElevator2--

		//check if any of the elevators have timed out and set the elevators active status to false
		if timeoutElevator0 <= 0 {
			temp[0] = false
		} else {
			temp[0] = true
		}
		if timeoutElevator1 <= 0 {
			temp[1] = false
		} else {
			temp[1] = true
		}
		if timeoutElevator2 <= 0 {
			temp[2] = false
		} else {
			temp[2] = true
		}

		for i := 0; i < 3; i++ {
			if temp[i] != prevTemp[i] {
				activeWatchdogs <- temp
				prevTemp = temp
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func WatchdogSendAlive(id int, watchdogTx chan int) {
	fmt.Println("Starting watchdog send alive")
	for {
		watchdogTx <- id
		time.Sleep(50 * time.Microsecond)
	}
}
