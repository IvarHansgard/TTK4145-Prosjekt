// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
)

func incrementing(x chan int, y chan bool) {
	//TODO: increment i 1000000 times
	for j := 0; j < 1000000; j++ {
		x <- 1
	}
	y <- true
}

func decrementing(x chan int, y chan bool) {
	//TODO: decrement i 1000000 times
	for k := 0; k < 1000000; k++ {
		x <- 1
	}
	y <- true
}

func server(inc, dec, result chan int, read chan bool) {
	i := 0
	for {
		select {
		case x := <-inc:
			i += x
		case x := <-dec:
			i -= x
		case <-read:
			result <- i
			return
		}
	}

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

	go incrementing(inc, done)
	go decrementing(dec, done)
	go server(inc, dec, result, read)

	<-done
	<-done
	read <- true

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.

	Println("The magic number is:", <-result)
}
