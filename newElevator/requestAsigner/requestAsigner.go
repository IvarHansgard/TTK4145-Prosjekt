package requestAsigner

import (
	"elevatorlib/elevator"
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
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func elevatorToHRAElevState(e elevator.Elevator) HRAElevState {
	var hra HRAElevState

	hra.Behavior = string(e.Behaviour)

	hra.Floor = e.Floor

	switch e.Dirn {
	case 0:
		hra.Direction = "up"
	case 1:
		hra.Direction = "down"
	case 2:
		hra.Direction = "stop"
	}
	for i := 0; i < 4; i++ {
		hra.CabRequests[i] = e.Requests[2][i]
	}

	return hra
}

func elevatorsToHRAInput(elevatorArray []elevator.Elevator) HRAInput {
	var input HRAInput
	input.States["one"] = elevatorToHRAElevState(elevatorArray[0])
	input.States["two"] = elevatorToHRAElevState(elevatorArray[1])
	input.States["three"] = elevatorToHRAElevState(elevatorArray[2])

	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			if elevatorArray[0].Requests[i][j] {
				input.HallRequests[i][j] = true
			}
			if elevatorArray[1].Requests[i][j] {
				input.HallRequests[i][j] = true
			}
			if elevatorArray[1].Requests[i][j] {
				input.HallRequests[i][j] = true
			}
		}
	}

	return input
}

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

func RequestAsigner(chActiveElevators chan []elevator.Elevator, masterState chan bool, chHallRequests chan HallRequests) {
	var oldActiveElevators []elevator.Elevator

	for {
		select {
		case elevators := <-chActiveElevators:
			if <-masterState {
				if checkIfNewRequests(elevators, oldActiveElevators) {
					input := elevatorsToHRAInput(elevators)

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

					oldActiveElevators = elevators
					chHallRequests <- *output
				}
			}
		//asign requests to elevators
		default:
			time.Sleep(10 * time.Second)

		}
	}
}
