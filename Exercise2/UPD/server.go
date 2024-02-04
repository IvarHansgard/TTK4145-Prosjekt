package main

//10.100.23.186

import (
	. "fmt"
	"log"
	"net"
	"runtime"
	"time"
)

func listenUDPport(port string, outsideIpAddress chan string) {
	recieveBuffer := make([]byte, 1024)

	udpAddress, err := net.ResolveUDPAddr("udp", port)

	udpConnection, err := net.ListenUDP("udp", udpAddress)

	numBytes, ipAddr, err := udpConnection.ReadFromUDP(recieveBuffer)
	Println(ipAddr.IP)
	if err != nil {
		log.Fatal(err, numBytes)
	}

	outsideIpAddress <- ipAddr.IP.String()
}

func sendUDP(ipAdress, port, message string) {
	udpAddress, err := net.ResolveUDPAddr("udp", ipAdress+port)
	udpConnection, err := net.ListenPacket("udp", ":0")

	numBytes, err := udpConnection.WriteTo([]byte(message), udpAddress)

	if err != nil {
		log.Fatal(err, numBytes)
	}

	Println("You sent: ", message, " Containing: ", numBytes, " bytes to IP: ", udpAddress.IP.String()+":", udpAddress.Port)

	udpConnection.Close()
}

/*
func broadcastdUDP(ipAdress, port, message string) {
}
*/

func readUDP(port string, buffer chan []byte, n chan int) {

	recieveBuffer := make([]byte, 1024)

	udpAddress, err := net.ResolveUDPAddr("udp", port)

	udpConnection, err := net.ListenUDP("udp", udpAddress)

	defer udpConnection.Close()

	numBytes, ipAddr, err := udpConnection.ReadFromUDP(recieveBuffer)

	if err != nil {
		log.Fatal(err, numBytes)
	}

	Println("you recieved: ", numBytes, "bytes from IP: ", ipAddr)

	buffer <- recieveBuffer
	n <- numBytes
}

func readUDP2(ip, port string) {
	pc, err := net.ListenPacket("udp", ip+port)
	if err != nil {
		panic(err)
	}
	defer pc.Close()
	buf := make([]byte, 1024)
	n, addr, err := pc.ReadFrom(buf)
	if err != nil {
		panic(err)
	}

	Println("you recieved: ", n, "bytes from IP: ", addr)
	Println("message was: ", buf[:n])
}

func sendTCP(ipAdress, port, message string) {
	tcpAdress, err := net.ResolveTCPAddr("tcp", ipAdress+port)
	tcpConnection, err := net.DialTCP("tcp", nil, tcpAdress)

	if err != nil {
		log.Fatal(err)
	}

	numBytes, err := tcpConnection.Write([]byte(message))
	tcpConnection.CloseWrite()

	if err != nil {
		log.Fatal(err, numBytes)
	}

}

func readTCP(ipAdress, port string) {
	recieveBuffer := make([]byte, 1024)

	tcpAdress, err := net.ResolveTCPAddr("tcp", ipAdress+port)
	tcpConnection, err := net.DialTCP("tcp", nil, tcpAdress)

	if err != nil {
		log.Fatal(err)
	}

	numBytes, err := tcpConnection.Read(recieveBuffer)
	tcpConnection.CloseRead()

	if err != nil {
		log.Fatal(err, numBytes)
	}
	Println(recieveBuffer[:numBytes])
}

func main() {

	runtime.GOMAXPROCS(2)
	//test listen udp to get server ip
	/*
		ipAddr := make(chan string)
		for {
			go listenUDPport(":30000", ipAddr)
			select {
			case x := <-ipAddr:
				Println(x)
			}
			time.Sleep(5 * time.Second)
		}
	*/

	//Test send/recieve UPD
	/*
		buffer := make(chan []byte, 1024)
		n := make(chan int)

		for {

			go sendUDP("10.100.23.191", ":20023", "HEI")

			go readUDP(":20023", buffer, n)

			select {
			case x := <-buffer:
				y := <-n
				print("Message recieved was: ")
				for i := 0; i < y; i++ {
					Printf("%c", x[i])
				}
				print("\n")
			}

			time.Sleep(5 * time.Second)
		}
	*/
	//Test send/recieve TCP
	for {
		//setup channels like udp ^^
		sendTCP("10.100.23.191", ":33546", "Connect to: 10.100.23.33")
		readTCP("10.100.23.191", ":33546")
		time.Sleep(5 * time.Second)
	}

}
