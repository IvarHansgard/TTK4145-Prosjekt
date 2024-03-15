package watchdog

import (
	"fmt"
	"time"
)

func WatchdogCheckAlive(elevatorSignal chan int, activeWatchdogs chan [3]bool) {
	fmt.Println("Starting watchdog check alive")
	elevator0Timer := time.NewTimer(5 * time.Second)
	elevator1Timer := time.NewTimer(5 * time.Second)
	elevator2Timer := time.NewTimer(5 * time.Second)

	for {
		select {
		case <-elevator0Timer.C:
			temp := <-activeWatchdogs
			temp[0] = false
			activeWatchdogs <- temp

			fmt.Println("elevator timer 1 went out")
		case <-elevator1Timer.C:
			temp := <-activeWatchdogs
			temp[1] = false
			activeWatchdogs <- temp

			fmt.Println("elevator timer 2 went out")
		case <-elevator2Timer.C:
			temp := <-activeWatchdogs
			temp[2] = false
			activeWatchdogs <- temp
			fmt.Println("elevator timer 3 went out")
		case id := <-elevatorSignal:
			switch id {
			case 0:
				temp := <-activeWatchdogs
				temp[0] = true
				activeWatchdogs <- temp
				elevator0Timer.Reset(5 * time.Second)
				fmt.Println("elevator timer 1 reset")
			case 1:
				temp := <-activeWatchdogs
				temp[1] = true
				activeWatchdogs <- temp
				fmt.Println("elevator timer 2 reset")
				elevator1Timer.Reset(5 * time.Second)
			case 2:
				temp := <-activeWatchdogs
				temp[2] = true
				activeWatchdogs <- temp
				fmt.Println("elevator timer 3 reset")
				elevator2Timer.Reset(5 * time.Second)
			default:
				break
			}
		}
	}
}

func WatchdogSendAlive(id int, watchdogTx chan int) {
	fmt.Println("Starting watchdog send alive")
	sendAliveTimer := time.NewTimer(1 * time.Millisecond)
	for {
		select {
		case <-sendAliveTimer.C:
			watchdogTx <- id
			sendAliveTimer.Reset(50 * time.Millisecond)

		}

	}
}
