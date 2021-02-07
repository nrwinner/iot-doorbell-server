package discover

import (
	"net"
)

const (
	MARCO = "marco"
	POLO  = "polo"
)

func StartDiscoveryServer() {
	pc, err := net.ListenPacket("udp4", ":9999")

	if err != nil {
		panic(err)
	}

	defer pc.Close()

	for {
		// create buffer to store received message, allocate a max of 1k space
		buffer := make([]byte, 1024)

		// read from the udp connection into the buffer
		packetLength, requester, err := pc.ReadFrom(buffer)

		if err != nil {
			panic(err)
		}

		// retrieve the packet as a string
		packet := string(buffer[:packetLength])

		if packet == MARCO {
			_, err := pc.WriteTo([]byte(POLO), requester)

			if err != nil {
				panic(err)
			}
		}
	}
}
