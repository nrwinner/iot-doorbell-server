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
					conn.WriteMessage(websocket.TextMessage, []byte("this is a test message"))
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				break
			}

			fmt.Println("recv:", message)
		}

		os.Exit(0)
	}()

	select {}
}
