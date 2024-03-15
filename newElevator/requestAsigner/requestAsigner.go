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

/*
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
*/
/*
func addHallRequests(button elevio.ButtonEvent) [2]int {
	var buttonArray [2]int
	buttonArray[0] = button.Floor
	buttonArray[1] = int(button.Button)
	return buttonArray
}
*/
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
*/ /*
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"*/
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

	input.States = inputStates
	input.HallRequests = hallRequest

	return input
}

func checkifNewHallRequest(choldHallRequests chan [4][2]bool, oldHallRequests, newHallRequests [4][2]bool) {

	fmt.Println("checking if new hall request")
	//fmt.Println("new requests:", newHallRequests)
	//fmt.Println("old requests:", oldHallRequests}
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

func setRunRequestAssigner(runHallRequestAssigner chan bool, state bool) {
	runHallRequestAssigner <- state
}

func RequestAsigner(chNewHallRequest chan elevio.ButtonEvent, chActiveElevators chan []elevator.Elevator, chMasterState chan bool,
	chClearedHallRequests chan elevio.ButtonEvent, hallRequestsTx chan HallRequests, chStopButtonPressed chan bool, chElevatorLost chan bool) {
	fmt.Println("Starting requestAsigner")

	choldHallRequests := make(chan [4][2]bool)

	hallRequests := [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
	oldHallRequests := [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}

	runHallRequestAssigner := make(chan bool)

	var elevatorStates []elevator.Elevator
	var masterState bool = true

	for {
		select {
		case <-chStopButtonPressed:
			go setRunRequestAssigner(runHallRequestAssigner, true)
		case temp := <-chMasterState:
			masterState = temp
			go setRunRequestAssigner(runHallRequestAssigner, true)

		case clearedHallRequest := <-chHallRequestClearedRx
:
			hallRequests[clearedHallRequest.Floor][int(clearedHallRequest.Button)] = false

			//elevio.SetButtonLamp(clearedHallRequest.Button,clearedHallRequest.Floor,false)

		case activeElevators := <-chElevatorStates:
			elevatorStates = activeElevators

		case button := <-chNewHallRequestRx:
			fmt.Println("Hall request recieved", button)
			hallRequests[button.Floor][int(button.Button)] = true
			go checkifNewHallRequest(choldHallRequests, oldHallRequests, hallRequests)

		case temp := <-choldHallRequests:
			oldHallRequests = temp
			//fmt.Println("old hall request set to", temp)
			go setRunRequestAssigner(runHallRequestAssigner, true)

		case newHallRequest := <-runHallRequestAssigner:
			//fmt.Println("newHallRequest is", newHallRequest)
			if masterState {
				//fmt.Println("new hall request")
				if newHallRequest {
					fmt.Println("Asigning requests to elevators")
					/*
						if len(oldHallRequests) == 0 {
							fmt.Println("No new requests")
							break
						}
					*/

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

					chAssignedHallRequestsTx<- *output
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
