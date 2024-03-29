package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"time"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:1234", nil)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		initPacket := CommandPacket{
			PacketType: "command",
			Command: "system/init",
			Args: map[string]string{
				"id": "doorbell01",
				"role": "camera",
			},
		}

		// send init packet
		conn.WriteJSON(initPacket)

		ticker := time.NewTicker(time.Second)
		quit := make(chan struct{})

		go func() {
			for {
				select {
				case <-ticker.C:
					conn.WriteJSON(CommandPacket{
						Id:         "doorbell01",
						PacketType: "command",
						Command:    "test/request",
						TargetDeviceId: "camera02",
					})
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		for {
			var packet CommandPacket
			var errorPacket ErrorPacket

			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				break
			}

			// unmarshal as a command and check packet type
			json.Unmarshal(message, &packet)

			if (packet.PacketType == "error") {
				// packet type is an error, unmarshal again as an error packet and handle
				json.Unmarshal(message, &errorPacket)
				fmt.Println("recv:", errorPacket)
			} else {
				fmt.Println("recv:", packet)
			}
		}

		os.Exit(0)
	}()

	select {}
}
