package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	f, err := os.Create("backup")
	check(err)

	defer f.Close()

	f.Sync()

	w := bufio.NewWriter(f)

	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()

	for i := 0; i < 10; i++ {
		s := fmt.Sprint(i, "\n")
		n, err := w.WriteString(s)
		w.Flush()
		print(n)
		check(err)
		time.Sleep(1 * time.Second)
	}

}
