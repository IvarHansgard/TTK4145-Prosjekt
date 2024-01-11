// Use `go run foo.go` to run your program
package main

import (
	. "fmt"
	"runtime"
)

func numberServer(inc, dec, result chan int, read chan bool) {
	i := 0
	for {
		select {
		case x := <-inc:
			i += x
		case x := <-dec:
			i += x
		case <-read:
			result <- i
			return
		}
	}
}

func increment(ch0 chan int, done chan bool) {
	for y := 0; y < 1000000; y++ {
		ch0 <- 1
	}
	done <- true
}

func decrement(ch0 chan int, done chan bool) {
	for y := 0; y < 1000000; y++ {
		ch0 <- -1
	}
	done <- true
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(3)
	// TODO: Spawn both functions as goroutines
	inc := make(chan int)
	dec := make(chan int)
	read := make(chan bool)
	result := make(chan int)
	done := make(chan bool, 2)

	go increment(inc, done)
	go decrement(dec, done)
	go numberServer(inc, dec, result, read)

	<-done
	<-done
	read <- true
	Println(<-result)

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort
	// We will do it properly with channels soon. For now: Sleep.

}
