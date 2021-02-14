package main

import (
	"doorbell-server/src/entities"
	"doorbell-server/src/modules/discover"
	"doorbell-server/src/modules/webrtc"
	"doorbell-server/src/modules/websocket"
	"fmt"
)

func main() {
	// start the websocket server and pass our controllers
	go websocket.WebSocketServer{}.StartServer(
		[]entities.Controller{
			// register the WebRTC module
			&webrtc.WebRTCController{},

			// the default case for testing purposes
			(&entities.DefaultController{}).New(
				func(client entities.Client, command entities.Command) {
					if command.Path == "test/request" {
						fmt.Println("Received:", command.Path)
						command.Responder.Respond(entities.Command{
							Path: "test/response",
							Args: map[string]string{"response": "response body"},
						})
					}
				},
				func(client entities.Client) { println("Connection", client.GetId()) },
				func(client entities.Client) { println("Disconnection", client.GetId()) },
			),
		},
	)

	// start the IP discovery server
	go discover.StartDiscoveryServer()

	fmt.Println("--- SERVER IS LISTENING ---")
	fmt.Println("Discovery - :9999")
	fmt.Println("WebSocketServer - :1234")

	// block forever
	select {}
}
