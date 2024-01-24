package main

import (
"fmt" 
"net"
)


func listnen(){

	recieveBuffer := make([]byte, 1024)

	udpConnection, err := net.ListenUDP("udp", udpAdress)

	numBytesRecieved, recieveAddr, err := udpConnection.ReadFrom(recieveBuffer)
	
	if(err.Error() != ""){
		fmt.Println("error")
	}
	
	fmt.Println(numBytesRecieved, recieveAddr.String())
	
}

func recieve(udpConnection net.UDPConn) []byte{
	recieveBuffer0, recieveBuffer1 := net.ReadMsgUDP(udpConnection)
	
	return recieveBuffer0
}

func send(writeBuffer []byte, address net.UDPAddr) int{
	answer, err = conn.WriteTo(writeBuffer, address)
	return answer 
}

func connectUDP(address, port string) net.UDPAddr{
	connection, err := net.ResolveUDPAddr("udp", address+port)

	if(err != ""){
		fmt.Println(err)
		break
	}
	
	return connection
}

func main(){
	listnen()
}