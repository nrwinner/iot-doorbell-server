package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Client struct {
	socket *websocket.Conn
	Role   string
	Id     string
}

func (c Client) ReadLoop(controller func(packet CommandPacket, client Client)) error {
	for {
		var packet CommandPacket
		err := c.socket.ReadJSON(&packet)

		if err != nil {
			// read error, assume disconnect
			return err
		} else {
			// pass message to socket controller
			controller(packet, c)
		}

	}
}

func (c Client) SendMessage(message string) {
	err := c.socket.WriteMessage(websocket.TextMessage, []byte(message))

	if err != nil {
		panic(err)
	}
}

func (c Client) SendCommand(command CommandPacket) {
	err := c.socket.WriteJSON(command)

	if err != nil {
		panic(err)
	}
}

func (c Client) SendError(error ErrorPacket) {
	err := c.socket.WriteJSON(error)

	if err != nil {
		panic(err)
	}
}

type WebSocketServer struct {
	connections []Client
}

func (s WebSocketServer) StartServer(controller Controller) {
	// set the default path to use our websocket handler
	http.HandleFunc("/", s.handleConnection(controller))
	err := http.ListenAndServe("localhost:1234", nil)

	if err != nil {
		// TODO:NickW better error handling here
		panic(err)
	}
}

func (s *WebSocketServer) handleConnection(controller Controller) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			// error upgrading HTTP request to WS
			panic(err)
		}

		defer conn.Close()

		// the first message on a new connection should be the initialization packet
		var packet InitPacket
		err = conn.ReadJSON(&packet)

		if err != nil || packet.PacketType != INIT_PACKET {
			socketErrorAndTerminate(conn, "could not decode init packet")
			return
		}

		// create new Client connection
		existingClient := s.locateClient(packet.Id)

		if existingClient != nil {
			// this should never happen, sanity check
			socketErrorAndTerminate(conn, "id already exists")
			return
		}

		client := Client{socket: conn, Role: packet.Role, Id: packet.Id}

		s.connections = append(
			s.connections,
			client,
		)

		controller.HandleConnect(client)

		// enter read loop and block
		readErr := client.ReadLoop(controller.HandleMessage)

		// at this point, readErr must be defined, but check for sanity
		if readErr != nil {
			s.disconnectClientById(client.Id, controller.HandleDisconnect)
		}
	}
}

func socketErrorAndTerminate(conn *websocket.Conn, message string) {
	_ = conn.WriteMessage(websocket.TextMessage, []byte("error - '"+message+"', closing connection"))
	conn.Close()
}

func (s *WebSocketServer) disconnectClientById(id string, disconnectHandler func(client Client)) {
	var newConnections []Client

	// remove this Client from list of connections
	for _, c := range s.connections {
		if c.Id != id {
			newConnections = append(newConnections, c)
		} else {
			c.socket.Close()
			disconnectHandler(c)
		}
	}

	s.connections = newConnections

}

func (s *WebSocketServer) locateClient(id string) *Client {
	for _, c := range s.connections {
		if c.Id == id {
			return &c
		}
	}

	// client not found
	return nil
}

func (s *WebSocketServer) locateAllClientsWithRole(role string) []*Client {
	var payload []*Client

	for _, c := range s.connections {
		if c.Role == role {
			payload = append(payload, &c)
		}
	}

	return payload
}
