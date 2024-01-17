package main

import (
"fmt" 
"net"
)

/*
var (sendingAddr = n.IPv4(0,0,0,0))
var (recieveAddr = n.IPv4(0,0,0,0))
*/

var udpAdress, err = net.ResolveUDPAddr("udp", "0.0.0.0:3000")


var recieveBuffer = make([]byte, 1024)

func listnen(){
	value, err := net.ListenUDP("udp", udpAdress)

	numBytesRecieved, recieveAddr, err := value.ReadFrom(recieveBuffer)
	
	if(err.Error() != ""){
		fmt.Println("error")
	}
	
	fmt.Println(numBytesRecieved, recieveAddr.String())
	
	}


func recieve(){

}

func send(){
	
}

func main(){
	listnen()
}