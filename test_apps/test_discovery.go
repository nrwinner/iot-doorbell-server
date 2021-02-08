package main

import (
	"fmt"
	"net"
)

const (
	BROADCAST_PORT = "9999"
	RECEIVE_PORT = "9998"
	BROADCAST      = "255.255.255.255"
	MARCO          = "marco"
	POLO           = "polo"
)

func main() {
	pc, err := net.ListenPacket("udp4", ":"+RECEIVE_PORT)

	if err != nil {
		panic(err)
	}

	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", BROADCAST+":"+BROADCAST_PORT)

	if err != nil {
		panic(err)
	}

	_, err = pc.WriteTo([]byte(MARCO), addr)

	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1024)

	packetLength, serverAddress, err := pc.ReadFrom(buffer)

	if err != nil {
		panic(err)
	}

	if string(buffer[:packetLength]) == POLO {
		serverAddress.String()
		fmt.Println("Address found!", serverAddress.String())
	}
}
