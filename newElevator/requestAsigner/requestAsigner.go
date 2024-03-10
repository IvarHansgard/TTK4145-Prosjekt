package requestAsigner

import (
	"elevatorlib/elevator"
	"elevatorlib/elevio"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HallRequests map[string][4][2]bool

type HRAElevState struct {
	Behavior    string  `json:"behaviour"`
	Floor       int     `json:"floor"`
	Direction   string  `json:"direction"`
	CabRequests [4]bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func elevatorToHRAElevState(e elevator.Elevator) HRAElevState {
	var hra HRAElevState
	fmt.Println("Converting elevator to HRAElevState:", e)

	if e.Behaviour == "" {
		e.Behaviour = elevator.EB_Disconnected
	}
	fmt.Println("Elevator behaviour:", string(e.Behaviour))
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

func elevatorsToHRAInput(hallRequest [][2]bool, elevatorArray []elevator.Elevator) HRAInput {
	var input HRAInput
	inputStates := make(map[string]HRAElevState)

	inputStates["one"] = elevatorToHRAElevState(elevatorArray[0])
	inputStates["two"] = elevatorToHRAElevState(elevatorArray[1])
	inputStates["three"] = elevatorToHRAElevState(elevatorArray[2])
	if inputStates["one"].Behavior == string(elevator.EB_Disconnected) {
		delete(inputStates, "one")
	}
	if inputStates["two"].Behavior == string(elevator.EB_Disconnected) {
		delete(inputStates, "two")
	}
	if inputStates["three"].Behavior == string(elevator.EB_Disconnected) {
		delete(inputStates, "three")
	}

	input.States = inputStates
	input.HallRequests = hallRequest

	return input
}

func compareHallRequests(oldHallRequests, newHallRequests [][2]bool) [][2]bool {
	fmt.Println("Comparing hall requests")
	fmt.Println("Old hall requests:", oldHallRequests)
	fmt.Println("New hall requests:", newHallRequests)

	for i := 0; i < len(newHallRequests); i++ {
		for j := 0; j < 2; j++ {
			if newHallRequests[i][j] != oldHallRequests[i][j] {
				oldHallRequests[i][j] = newHallRequests[i][j]
			}
		}
	}
	return oldHallRequests
}

/*
func checkIfNewRequests(elevators, oldActiveElevators []elevator.Elevator) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 2; k++ {
				if elevators[i].Requests[j][k] != oldActiveElevators[i].Requests[j][k] {
					return true
				}
			}
		}
	}
	return false
}
*/
/*
	func getHallRequests(elevators []elevator.Elevator) [4][2]bool {
		var hallRequests [4][2]bool
		for i := 0; i < 3; i++ {
			for j := 0; j < 4; j++ {
				for k := 0; k < 2; k++ {
					if elevators[i].Requests[j][k] {
						hallRequests[j][k] = true
					}
				}
			}					elevatorTx <- localElevator

		}
		return hallRequests
	}
*/
func RequestAsigner(chActiveElevators chan []elevator.Elevator, masterState bool, localHallRequestsRx chan [][2]bool, hallRequestsTx chan HallRequests, chHallRequestCleared chan [2]int) {
	fmt.Println("Starting requestAsigner")
	oldHallRequests := [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}

	for {
		select {
		case localHallrequests := <-localHallRequestsRx:
			for i := 0; i < len(localHallrequests); i++ {
				for j := 0; j < 2; j++ {
					if localHallrequests[i][j] {
						if !oldHallRequests[i][j] {
							oldHallRequests[i][j] = true
						}
					}
				}
			}
		case hallRequestCleared := <-chHallRequestCleared:
			oldHallRequests[hallRequestCleared[0]][hallRequestCleared[1]] = false

		case elevators := <-chActiveElevators:
			if masterState {
				fmt.Println("Asigning requests to elevators")
				/*
					if len(oldHallRequests) == 0 {
						fmt.Println("No new requests")
						break
					}
				*/

				input := elevatorsToHRAInput(oldHallRequests, elevators)

				hraExecutable := ""
				switch runtime.GOOS {
				case "linux":
					hraExecutable = "hall_request_assigner"
				case "windows":
					hraExecutable = "hall_request_assigner.exe"
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
				fmt.Println("Hall requests assigned: ", *output)
				hallRequestsTx <- *output
			}
			//asign requests to elevators
		default:
			time.Sleep(10 * time.Second)
		}
	}
}
