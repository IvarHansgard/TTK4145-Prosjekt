package main

import (
	"fmt"
	"time"
)

func producer(x chan int) {

	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("[producer]: pushing %d\n", i)
		// TODO: push real value to buffer
		x <- i
	}

}

func consumer(x chan int) {

	time.Sleep(1 * time.Second)
	for { //TODO: get real value from buffer
		i := <-x
		fmt.Printf("[consumer]: %d\n", i)
		time.Sleep(50 * time.Millisecond)
	}

}

func main() {

	// TODO: make a bounded buffer
	buffer := make(chan int, 5)

	go consumer(buffer)
	go producer(buffer)

	select {}
}
//hei jeg heter ulrikke