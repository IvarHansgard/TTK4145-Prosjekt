package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	number := 0
	/*
		f, err := os.Create("./backup.txt")

		if err != nil {
			panic(err)
		}

		f.Close()
	*/
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
	if err != nil {
		panic(err)
	}

	for i := 0; i < 6; i++ {
		fmt.Println(number)
		number = number + 1
		time.Sleep(1 * time.Second)
	}
}
