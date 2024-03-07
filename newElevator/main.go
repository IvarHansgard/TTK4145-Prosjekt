package make

import (
	elevator "ElevatorLib/Elevator"
	"ElevatorLib/elevator"
	"ElevatorLib/network/bcast"
	request_asigner "ElevatorLib/requestAsigner"
	"ElevatorLib/watchdog"
	"flag"
)

func checkMaster(id int, chActiveElevators chan []bool) bool {
	if chActiveElevators[id-1] == true && id != 0 {
		return false
	} else {
		return true
	}
}

func main() {
	var id int = 0
	flag.IntVar(&id, "id", "", "id of this elevator")
	flag.Parse()

	master := make(chan boot)
	if id == 0 {
		masterState <- true
	}

	elevatorTx = make(chan elevator.Elevator)
	elevatorRx = make(chan elevator.Elevator)
	activeElevators := [3]elevator.Elevator
	chActiveElevators = make(chan []elevator.Elevator)

	watchdogTx = make(chan int)
	watchdogRx = make(chan int)
	chActiveWatchdogs = make(chan []bool)

	localElevator = elevator.Elevator_init(id)
	activeElevators[id] = localElevator
	chLocalElevator = make(chan elevator.Elevator)
	chLocalElevator <- localElevator

	go bcast.Transmitter(2000, elevatorTx)
	go bcast.Receiver(2001, elevatorRx)

	go bcast.Transmitter(3000, watchdogTx)
	go bcast.Receiver(3001, watchdogRx)

	go watchdog.Watchdog_checkAlive(watchdogRx, chActiveWatchdogs, 10)
	go watchdog.Watchdog_sendAlive(id, watchdogTx)

	go elvator.RunElevator(chLocalElevator)

	go request_asigner.Request_asigner(chActiveElevators, elevatorTx) //jobbe med den her

	for {
		select {
		case elevator := <-elevatorRx:
			activeElevators[elevator.Id] = elevator
			if masterState == true {
				chActiveElevators <- activeElevators
			}

		case activeWatchdogs := <-chActiveWatchdogs:
			masterState <- checkMaster(id, activeElevators)

		case masterState := <-masterState:
			if masterState == true {
				changeToMaster()
			} else {
				changeToSlave()
			}
		}
	}
}
