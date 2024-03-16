Running the elevator
==========================================
To run the elevator type without compiling type `go run main.go <optional arguments>` or to compile type `go build main.go` 

After compiling you can run with `./main.exe <optional arguments>` on windows or `./main <optional arguments>` on linux 

Arguments avalible
==========================================
Id of the elevator (Default = "0"): `--id`

Port to the elevator server you want to use (Default = 15657): `--port`

Number of elevators in the system (Default = 3): `--elevators`

Example
==========================================
Running with "3" elevators on port "15651" and an id of "1"

Without compiling: `go run main.go --elevators 3 --port 15651 --id 1`

After compiling:

Windows: `./main.exe --elevators 3 --port 15651 --id 1`

Linux: `./main --elevators 3 --port 15651 --id 1`