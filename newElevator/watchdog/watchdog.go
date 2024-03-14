package watchdog

import (
	"fmt"
	"time"
)

func WatchdogCheckAlive(elevatorSignal chan int, activeWatchdogs chan [3]bool, timeout int) {
	fmt.Println("Starting watchdog check alive")

	elevator0Timer := time.NewTimer(time.Duration(timeout) * time.Second)
	elevator1Timer := time.NewTimer(time.Duration(timeout) * time.Second)
	elevator2Timer := time.NewTimer(time.Duration(timeout) * time.Second)

	for {
		select {
		case <-elevator0Timer.C:
			temp := <-activeWatchdogs
			temp[0] = false
			activeWatchdogs <- temp
		case <-elevator1Timer.C:
			temp := <-activeWatchdogs
			temp[1] = false
			activeWatchdogs <- temp
		case <-elevator2Timer.C:
			temp := <-activeWatchdogs
			temp[2] = false
			activeWatchdogs <- temp
		case id := <-elevatorSignal:
			switch id {
			case 0:
				elevator0Timer.Reset(time.Duration(timeout) * time.Second)
			case 1:
				elevator1Timer.Reset(time.Duration(timeout) * time.Second)
			case 2:
				elevator2Timer.Reset(time.Duration(timeout) * time.Second)
			}
		}
	}
}

func WatchdogSendAlive(id int, watchdogTx chan int) {
	fmt.Println("Starting watchdog send alive")
	for {
		watchdogTx <- id
		time.Sleep(50 * time.Millisecond)
	}
}
