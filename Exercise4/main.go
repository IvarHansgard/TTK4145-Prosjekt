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

/*
func mainA() {

	number := 0
	/*
		f, err := os.Create("./backup.txt")

		if err != nil {
			panic(err)
		}

		f.Close()

	dat, err := os.Open("./backup.txt")

	if err != nil {
		panic(err)
	}

	//fileReader := bufio.NewReader(dat)
	//fileWriter := bufio.NewWriter(dat)
	//file := bufio.NewReadWriter(fileReader, fileWriter)
	n, err := dat.Write([]byte("string"))
	//n, err = fileWriter.WriteString("0string(number)\n")
	fmt.Println(n)
	/*
		s, err := fileReader.ReadBytes('\n')
		number, err = strconv.Atoi(string(s))
*/
/*
		if err != nil {
			panic(err)
		}

	for i := 0; i < 6; i++ {
		fmt.Println(number)
		number = number + 1
		time.Sleep(1 * time.Second)
	}
}
*/
// funksjon for A
func mainA() {

	number := 0

	for i := 0; i < 6; i++ {
		fmt.Println(number)
		number = number + 1
		time.Sleep(1 * time.Second)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	number := 0

	if _, err := os.Stat("backup"); err == nil {

		for {
			info, err := os.Stat("backup")
			check(err)

			t := time.Now().Second()

			if t == info.ModTime().Second()+3 {
				break
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
			exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
		}
	} else if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create("backup")
		check(err)

		defer f.Close()

		f.Sync()

		exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	}
	f, err := os.Open("backup")
	check(err)

	w := bufio.NewWriter(f)

	for i := number; i < i+10; i++ {
		s := fmt.Sprint(i, "\n")
		n, err := w.WriteString(s)
		w.Flush()
		print(n)
		check(err)
		time.Sleep(1 * time.Second)
	}

}
