package main

import (
	"doorbell-server/src/modules/websocket"
	"fmt"
)

func main() {
	// start the websocket server and pass our controller function
	go websocket.WebSocketServer{}.StartServer(SocketMessageController)

	// block forever
	select {}
}

func SocketMessageController(message string) {
	// TODO:NickW move this and implement
	fmt.Println("Controller", message)
}
