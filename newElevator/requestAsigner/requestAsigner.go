package requestAsigner

import (
	"elevatorlib/elevator"
	"elevatorlib/elevio"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type HallRequests map[string][4][2]bool
type ElevatorMap map[string]elevator.Elevator

type HRAElevState struct {
	Behavior    string  `json:"behaviour"`
	Floor       int     `json:"floor"`
	Direction   string  `json:"direction"`
	CabRequests [4]bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [4][2]bool              `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func elevatorToHRAElevState(e elevator.Elevator) HRAElevState {
	var hra HRAElevState

	if e.Behaviour == "" {
		e.Behaviour = elevator.EB_Disconnected
	}
	hra.Behavior = string(e.Behaviour)
	hra.Floor = e.Floor

	switch e.Dirn {
	case elevio.MD_Up:
		hra.Direction = "up"
	case elevio.MD_Down:
		hra.Direction = "down"
	case elevio.MD_Stop:
		hra.Direction = "stop"
	}
	for i := 0; i < 4; i++ {
		hra.CabRequests[i] = e.Requests[i][2]
	}

	return hra
}

func elevatorsToHRAInput(hallRequest [4][2]bool, elevatorArray ElevatorMap) HRAInput {
	var input HRAInput
	inputStates := make(map[string]HRAElevState)

	for i := 0; i < len(elevatorArray); i++ {
		if elevatorArray[strconv.Itoa(i)].Behaviour != elevator.EB_Disconnected {
			inputStates[elevatorArray[strconv.Itoa(i)].StrID] = elevatorToHRAElevState(elevatorArray[strconv.Itoa(i)])
		}
	}
	input.States = inputStates
	input.HallRequests = hallRequest

	return input
}

func checkifNewHallRequest(oldHallRequests, newHallRequests [4][2]bool) bool {
	//debug
	//fmt.Println("checking if new hall request")
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			if !oldHallRequests[i][j] == newHallRequests[i][j] {
				return true
			}

		}

	}
	return false
}

func setRunRequestAssigner(runHallRequestAssigner chan bool, state bool) {
	runHallRequestAssigner <- state
}

func RequestAsigner(
	chNewHallRequestRx chan elevio.ButtonEvent,
	chElevatorRx chan elevator.Elevator,
	chMasterState chan bool,
	chHallRequestClearedRx chan elevio.ButtonEvent,
	chAssignedHallRequestsTx chan HallRequests,
	chStopButtonPressed chan bool,
	chSendHallRequestsToMasterTx chan [4][2]bool,
	chSendHallRequestsToMasterRx chan [4][2]bool,
	chSendElevatorStatesToMasterTx chan ElevatorMap,
	chSendElevatorStatesToMasterRx chan ElevatorMap,
	chElevatorLost chan string,
	numElevators int) {

	fmt.Println("Starting requestAsigner")

	hallRequests := [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
	oldHallRequests := [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}

	runHallRequestAssigner := make(chan bool)
	elevatorStates := make(map[string]elevator.Elevator)
	for i := 0; i < numElevators; i++ {
		elevatorStates[strconv.Itoa(i)] = elevator.Elevator_init_disconnected(i, strconv.Itoa(i))
	}
	var masterState bool = true

	for {
		select {
		case <-chStopButtonPressed:
			go setRunRequestAssigner(runHallRequestAssigner, true)

		case temp := <-chMasterState:
			if masterState && !temp {
				fmt.Println("sending data to master")
				chSendHallRequestsToMasterTx <- hallRequests
				chSendElevatorStatesToMasterTx <- elevatorStates
				masterState = temp
			}

		case hallRequesFromPrevMaster := <-chSendHallRequestsToMasterRx:
			if masterState {
				hallRequests = hallRequesFromPrevMaster
			}

		case clearedHallRequest := <-chHallRequestClearedRx:
			//debug
			//fmt.Println("clearing button", clearedHallRequest.Floor, int(clearedHallRequest.Button))
			//oldHallRequests[clearedHallRequest.Floor][int(clearedHallRequest.Button)] = false
			hallRequests[clearedHallRequest.Floor][int(clearedHallRequest.Button)] = false

		case elevator := <-chElevatorRx:
			elevatorStates[elevator.StrID] = elevator

		case elevatorStatesFromPrevMaster := <-chSendElevatorStatesToMasterRx:
			elevatorStates = elevatorStatesFromPrevMaster

		case lostElevatorID := <-chElevatorLost:
			tempElev := elevatorStates[lostElevatorID]
			tempElev.Behaviour = elevator.EB_Disconnected
			elevatorStates[lostElevatorID] = tempElev
			if masterState {
				go setRunRequestAssigner(runHallRequestAssigner, true)
			}

		case button := <-chNewHallRequestRx:
			hallRequests[button.Floor][int(button.Button)] = true
			if checkifNewHallRequest(oldHallRequests, hallRequests) {
				oldHallRequests = hallRequests
				go setRunRequestAssigner(runHallRequestAssigner, true)
			}

		case run := <-runHallRequestAssigner:

			if masterState {

				if run {
					fmt.Println("Asigning requests to elevators")

					input := elevatorsToHRAInput(hallRequests, elevatorStates)

					hraExecutable := ""
					switch runtime.GOOS {
					case "linux":
						hraExecutable = "hall_request_assigner"
					case "windows":
						hraExecutable = "./requestAsigner/hall_request_assigner.exe"
					default:
						panic("OS not supported")
					}

					jsonBytes, err := json.Marshal(input)
					if err != nil {
						fmt.Println("json.Marshal error: ", err)
						return
					}

					ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
					if err != nil {
						fmt.Println("exec.Command error: ", err)
						fmt.Println(string(ret))
						return
					}

					output := new(map[string][4][2]bool)
					err = json.Unmarshal(ret, &output)
					if err != nil {
						fmt.Println("json.Unmarshal error: ", err)
						return
					}

					chAssignedHallRequestsTx <- *output
					//debug
					fmt.Println("Hall requests assigned: ", *output)
				}
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
