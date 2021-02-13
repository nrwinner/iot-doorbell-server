package websocket

import (
	"doorbell-server/src/entities"
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	Id     string
	Role   string
	socket *websocket.Conn
	mutex sync.Mutex
}

func (c *Client) GetId() string {
	return c.Id
}

func (c *Client) GetRole() string {
	return c.Role
}

func (c *Client) SendMessage(message string) {
	err := c.socket.WriteMessage(websocket.TextMessage, []byte(message))

	if err != nil {
		panic(err)
	}
}

func (c *Client) SendCommand(command entities.Command) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	packet := PacketFromCommand(command)
	err := c.socket.WriteJSON(packet)

	if err != nil {
		panic(err)
	}
}

func (c *Client) SendError(error string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	packet := ErrorPacket{
		Packet: Packet{PacketType: ERROR_PACKET},
		Error:  error,
	}
	err := c.socket.WriteJSON(packet)

	if err != nil {
		panic(err)
	}
}
