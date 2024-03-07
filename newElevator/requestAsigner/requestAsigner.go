package request_asigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
    "ElevatorLib/elevator"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

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


func elevatorToHRA(es <- chan Elevator) HRAElevState{
    var hra HRAElevState
    hra.Behavior = es.behaviour
    hra.CabRequests = es.requests[3]
    hra.Direction = es.dirn
    hra.Floor = es.floor
    return hra
}

func Elevator_algo(hallRequests <- chan [4][2]int, elevator1, elevator2, elevator3 Elevator) out map[string][][2]bool{
    hra1  := elevatorToHRA(elevator1)
    hra2  := elevatorToHRA(elevator2)
    hra3  := elevatorToHRA(elevator3)

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	input := HRAInput{
		HallRequests: [][2]bool{{hallRequests[0][0], hallRequests[0][1]}, {hallRequests[1][0], hallRequests[1][1]}, {hallRequests[2][0], hallRequests[2][1]}, {hallRequests[3][0], hallRequests[3][1]}},
		States: map[string]HRAElevState{
			"one": hra1,
			"two": hra2,
			"three": hra3,
		},
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return
	}

	ret, err := exec.Command("../hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
    return output
}
