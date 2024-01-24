package main

//10.100.23.186

import (
	. "fmt"
	"log"
	"net"
)

func listenUDPport(port string) string {
	recieveBuffer := make([]byte, 1024)

	udpAddress, err := net.ResolveUDPAddr("udp", ":"+port)

	udpConnection, err := net.ListenUDP("udp", udpAddress)

	numBytes, ipAddr, err := udpConnection.ReadFromUDP(recieveBuffer)

	if err != nil {
		log.Fatal(err, numBytes)
	}

	return ipAddr.IP.String()
}

func sendUDP() {
	recieveBuffer := make([]byte, 1024)

	udpAddress, err := net.ResolveUDPAddr("udp", "10.100.23.186:20023")
	udpConnection, err := net.ListenPacket("udp", ":0")

	numBytes, err := udpConnection.WriteTo([]byte("HEI"), udpAddress)

	if err != nil {
		log.Fatal(err, numBytes)
	}
	Print(numBytes)
}

func main() {

}
