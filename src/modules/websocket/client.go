package websocket

import (
	"doorbell-server/src/entities"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id     string
	Role   string
	socket *websocket.Conn
}

func (c Client) GetId() string {
	return c.Id
}

func (c Client) GetRole() string {
	return c.Role
}

func (c Client) SendMessage(message string) {
	err := c.socket.WriteMessage(websocket.TextMessage, []byte(message))

	if err != nil {
		panic(err)
	}
}

func (c Client) SendCommand(command entities.Command) {
	packet := PacketFromCommand(command)
	packet.Id = c.Id
	err := c.socket.WriteJSON(packet)

	if err != nil {
		panic(err)
	}
}

func (c Client) SendError(error string) {
	packet := ErrorPacket{
		Packet: Packet{Id: c.Id, PacketType: ERROR_PACKET},
		Error:  error,
	}
	err := c.socket.WriteJSON(packet)

	if err != nil {
		panic(err)
	}
}
