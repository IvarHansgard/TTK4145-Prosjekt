package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	number := 0

	if _, err := os.Stat("backup"); err == nil {
		// path/to/whatever exists

		for {
			info, err := os.Stat("backup")
			check(err)

			t := time.Now().Second()

			if t == info.ModTime().Second()+3 {
				break
			}
		}

		var lastline string

		f, err := os.Open("backup")
		check(err)
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			lastline = scanner.Text()
		}

		n, err := strconv.Atoi(lastline)
		check(err)
		number = n
		//exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		f, err := os.OpenFile("backup", os.O_CREATE, 0600)
		check(err)

		defer f.Close()

		f.Sync()

		//exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()

	} else {
		fmt.Println("ERROR!!")
	}

	f, err := os.OpenFile("backup", os.O_WRONLY|os.O_APPEND, 0600)
	check(err)

	w := bufio.NewWriter(f)

	//exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go").Run()
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()

	for i := number; i < number+10; i++ {
		s := fmt.Sprint(i, "\n")
		n, err := w.WriteString(s)
		w.Flush()
		fmt.Println("wrote: ", n, " bytes", " number: ", i)
		check(err)
		time.Sleep(1 * time.Second)
	}

}
