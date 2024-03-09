package make

import (
	"elevatorlib/elevator"
	"elevatorlib/elevator/runElevator"
	"elevatorlib/network/bcast"
	"elevatorlib/requestAsigner"
	"elevatorlib/watchdog"
	"flag"
)

func checkMaster(id int, activeElevators []bool) bool {
	if activeElevators[id-1] == true && id != 0 {
		return false
	} else {
		return true
	}
}

type hallRequests map[string][][2]int

func main() {
	var id int
	flag.IntVar(&id, "id", 0, "id of this elevator")
	flag.Parse()

	masterState := make(chan bool)
	if id == 0 {
		masterState <- true
	}

	elevatorTx := make(chan elevator.Elevator)
	elevatorRx := make(chan elevator.Elevator)
	activeElevators := make([]elevator.Elevator, 3)
	chActiveElevators := make(chan []elevator.Elevator)

	hallRequestsTx := make(chan requestAsigner.HallRequests)
	hallRequestsRx := make(chan requestAsigner.HallRequests)

	watchdogTx := make(chan int)
	watchdogRx := make(chan int)
	chActiveWatchdogs := make(chan []bool)

	go bcast.Transmitter(2000, elevatorTx)
	go bcast.Receiver(2001, elevatorRx)

	go bcast.Transmitter(3000, hallRequestsTx)
	go bcast.Receiver(3001, hallRequestsRx)

	go bcast.Transmitter(4000, watchdogTx)
	go bcast.Receiver(4001, watchdogRx)

	go watchdog.Watchdog_checkAlive(watchdogRx, chActiveWatchdogs, 10)
	go watchdog.Watchdog_sendAlive(id, watchdogTx)

	go runElevator.RunLocalElevator(chActiveElevators, elevatorTx, hallRequestsRx, id)

	go requestAsigner.RequestAsigner(chActiveElevators, masterState, hallRequestsTx) //jobbe med den her

	for {
		select {
		case elevator := <-elevatorRx:
			activeElevators[elevator.Id] = elevator
			chActiveElevators <- activeElevators

		case activeWatchdogs := <-chActiveWatchdogs:
			masterState <- checkMaster(id, activeWatchdogs)
		}
	}
}
