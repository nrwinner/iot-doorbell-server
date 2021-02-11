package main

import (
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
		initPacket := make(map[string]string)

		initPacket["PacketType"] = "init"
		initPacket["Role"] = "camera"
		initPacket["Id"] = "doorbell01"

		if err != nil {
			panic(err)
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
						Args:       []string{},
					})
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		for {
			var packet CommandPacket
			err := conn.ReadJSON(&packet)
			if err != nil {
				fmt.Println("read:", err)
				break
			}

			fmt.Println("recv:", packet)
		}

		os.Exit(0)
	}()

	select {}
}
