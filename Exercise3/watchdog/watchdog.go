package watchdog

import (
	"time"
)

func watchdog(elevatorSignal, timeoutSignal chan int, timeout int) {
	initialTimeout := timeout
	for {
		select {
		case <-elevatorSignal:
			timeout = initialTimeout
		}

		timeout--

		if timeout == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	timeoutSignal <- 1
}

func watchdog_sendAlive(id int, tx chan int) {
	tx <- id
}
