package websocket

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

const (
	CAMERA = 0
	VIEWER = 1
	UPDATE_MANAGER = 2
)

type Client struct {
	socket *websocket.Conn
	Role   int
	Id     string
}

func (c Client) ReadLoop(controller func(message string)) error {
	for {
		_, message, err := c.socket.ReadMessage()

		if err != nil {
			// read error, assume disconnect
			return err
		} else {
			// pass message to socket controller
			controller(string(message[:]))
			c.SendMessage("ack")
		}

	}
}

func (c Client) SendMessage(message string) {
	err := c.socket.WriteMessage(websocket.TextMessage, []byte(message))

	if err != nil {
		// TODO:NickW better error handling here
		panic(err)
	}
}

type WebSocketServer struct {
	connections []Client
}

func (s WebSocketServer) StartServer(controller func(message string)) {
	// set the default path to use our websocket handler
	http.HandleFunc("/", s.handleConnection(controller))
	err := http.ListenAndServe("localhost:1234", nil)

	if err != nil {
		// TODO:NickW better error handling here
		panic(err)
	}
}

func (s *WebSocketServer) handleConnection(controller func(message string)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(len(s.connections))
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			// error upgrading HTTP request to WS
			panic(err)
		}

		defer conn.Close()

		// the first message on a new connection should be the initialization packet
		_, initPacket, err := conn.ReadMessage()

		if err != nil {
			panic(err)
		}

		// parse the init packet
		role, id, parseError := parseInitPacket(string(initPacket[:]))

		if parseError != nil {
			// TODO:NickW handle init-packet parse error better
			fmt.Println("Parse error with packet", string(initPacket[:]))
		}

		// create new Client connection
		client := Client{socket: conn, Role: role, Id: id}
		existingClient := s.locateClient(id)

		if existingClient != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Connection with id " + id + " already exists"))
			conn.Close()
			return
		}

		s.connections = append(
			s.connections,
			client,
		)

		fmt.Println("Connections:", len(s.connections))

		// enter read loop and block
		readErr := client.ReadLoop(controller)

		// at this point, readErr must be defined, but check for sanity
		if readErr != nil {
			s.disconnectClientById(client.Id)
			fmt.Println("Connections:", len(s.connections))
		}
	}
}

func parseInitPacket(initPacket string) (role int, id string, err error) {
	if initPacket[:4] != "init" {
		panic("No initialization packet detected")
	}

	// packet format ["init"][2-digit role][name]
	// example init00camera-back
	role, roleErr := strconv.Atoi(initPacket[4:6])
	id = initPacket[6:]

	if roleErr != nil || id == "" {
		// init packet could not be parsed
		return -1, "", errors.New("improperly-formatted initialization packet")
	}

	return role, id, nil
}

func (s *WebSocketServer) disconnectClientById(id string) {
	var newConnections []Client

	// remove this Client from list of connections
	for _, c := range s.connections {
		if c.Id != id {
			newConnections = append(newConnections, c)
		} else {
			c.socket.Close()
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

func (s *WebSocketServer) locateAllClientsWithRole(role int) []*Client {
	var payload []*Client

	for _, c := range s.connections {
		if c.Role == role {
			payload = append(payload, &c)
		}
	}

	return payload
}
