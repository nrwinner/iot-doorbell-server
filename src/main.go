package main

import (
	"doorbell-server/src/modules/discover"
	"doorbell-server/src/modules/websocket"
)

func main() {
	// start the websocket server and pass our controller function
	go websocket.WebSocketServer{}.StartServer(
		websocket.Controller{},
	)

	// start the IP discovery server
	go discover.StartDiscoveryServer()

	// block forever
	select {}
}
