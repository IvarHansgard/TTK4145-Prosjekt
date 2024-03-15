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
	HallRequests [4][2]bool              `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func elevatorToHRAElevState(e elevator.Elevator) HRAElevState {
	var hra HRAElevState
	//fmt.Println("Converting elevator to HRAElevState:", e)

	if e.Behaviour == "" {
		e.Behaviour = elevator.EB_Disconnected
	}
	//fmt.Println("Elevator behaviour:", string(e.Behaviour))
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

func elevatorsToHRAInput(hallRequest [4][2]bool, elevatorArray []elevator.Elevator) HRAInput {
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
	fmt.Println("elevator one state", inputStates["one"].Behavior)
	fmt.Println("elevator two state", inputStates["two"].Behavior)
	fmt.Println("elevator three state", inputStates["three"].Behavior)
	input.States = inputStates
	input.HallRequests = hallRequest

	return input
}

func checkifNewHallRequest(choldHallRequests chan [4][2]bool, oldHallRequests, newHallRequests [4][2]bool) {

	fmt.Println("checking if new hall request")
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			if !oldHallRequests[i][j] == newHallRequests[i][j] {
				fmt.Println("true")
				choldHallRequests <- newHallRequests
				return
			}

		}

	}
}

func setRunRequestAssigner(isNewHallRequest chan bool, state bool) {
	isNewHallRequest <- state
}

func RequestAsigner(chNewHallRequestRx chan elevio.ButtonEvent, chElevatorStates chan []elevator.Elevator, chRequestAssignerMasterState chan bool,
	chHallRequestClearedRx chan elevio.ButtonEvent, chAssignedHallRequestsTx chan HallRequests, chStopButtonPressed chan bool) {
	fmt.Println("Starting requestAsigner")

	chOldHallRequests := make(chan [4][2]bool)

	hallRequests := [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
	oldHallRequests := [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}

	runHallRequestAssigner := make(chan bool)

	var elevatorStates []elevator.Elevator
	var masterState bool = true

	for {
		select {
		case <-chStopButtonPressed:
			go setRunRequestAssigner(runHallRequestAssigner, true)
		case temp := <-chRequestAssignerMasterState:
			masterState = temp
			go setRunRequestAssigner(runHallRequestAssigner, true)

		case clearedHallRequest := <-chHallRequestClearedRx:
			hallRequests[clearedHallRequest.Floor][int(clearedHallRequest.Button)] = false

		case activeElevators := <-chElevatorStates:
			elevatorStates = activeElevators

		case button := <-chNewHallRequestRx:
			fmt.Println("Hall request recieved", button)
			hallRequests[button.Floor][int(button.Button)] = true
			go checkifNewHallRequest(chOldHallRequests, oldHallRequests, hallRequests)

		case temp := <-chOldHallRequests:
			oldHallRequests = temp
			go setRunRequestAssigner(runHallRequestAssigner, true)

		case newHallRequest := <-runHallRequestAssigner:

			if masterState {

				if newHallRequest {
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
					//fmt.Println("Hall requests assigned: ", *output)
					//fmt.Println("old", oldHallRequests)
					//fmt.Println("new", HallRequests)
				}
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
