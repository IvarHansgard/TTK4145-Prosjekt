Running the elevator
==========================================
Before running the program go into main.go and uncomment / comment line `100` or `102` depending on your operating system and save the file.

For windows the lines `98` to `102` should look like this:
```
//start new elevator
//Windows
	go exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go", "--id "+id, "--port "+strconv.Itoa(port), "--elevators "+strconv.Itoa(numElevators)).Run()
//Linux
	//exec.Command("gnome-terminal", "--", "go", "run", "main.go", "--id "+id, "--port "+strconv.Itoa(port), "--elevators "+strconv.Itoa(numElevators)).Run()
```
and for linux like this:
```
//start new elevator
//Windows
    //go exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go", "--id "+id, "--port "+strconv.Itoa(port), "--elevators "+strconv.Itoa(numElevators)).Run()
//Linux
    exec.Command("gnome-terminal", "--", "go", "run", "main.go", "--id "+id, "--port "+strconv.Itoa(port), "--elevators "+strconv.Itoa(numElevators)).Run()
```

To run the elevator type without compiling type `go run main.go <optional arguments>` or to compile type `go build main.go` 

After compiling you can run with `./main.exe <optional arguments>` on windows or `./main <optional arguments>` on linux 


Arguments avalible
==========================================
Id of the elevator (Default = "0"): `--id`

Port to the elevator server you want to use (Default = 15657): `--port`

Number of elevators in the system (Default = 3): `--elevators`

Timeout for crash detection in seconds, i.e the time it will take before your backup will take over after detecting a crash of the main program  (Default = 5): `--timeout`


Example
==========================================
Running with "3" elevators on port "15651" and an id of "1"

Without compiling: `go run main.go --elevators 3 --port 15651 --id 1`

After compiling:

Windows: `./main.exe --elevators 3 --port 15651 --id 1`

Linux: `./main --elevators 3 --port 15651 --id 1`