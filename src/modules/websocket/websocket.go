package websocket

import (
	"doorbell-server/src/entities"
	"fmt"
	"net/http"

	ws "github.com/gorilla/websocket"
)

type WebSocketServer struct {
	connections []Client
}

const (
	initPacket = "system/init"
)

func (s WebSocketServer) StartServer(controllers []entities.Controller) {
	// set the default path to use our websocket handler
	http.HandleFunc("/", s.handleConnection(controllers))
	err := http.ListenAndServe(":1234", nil)

	if err != nil {
		// TODO(nrwinner) better error handling here
		panic(err)
	}
}

func (s *WebSocketServer) handleConnection(controllers []entities.Controller) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// FIXME(nrwinner) this is extremely insecure, lets not blindly trust all origins
		upgrader := ws.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			// error upgrading HTTP request to WS
			fmt.Println("ORIGIN WOOHOO", r.Header.Get("origin"))
			panic(err)
		}

		defer conn.Close()

		// the first message on a new connection should be the initialization packet
		var packet CommandPacket
		err = conn.ReadJSON(&packet)

		id, hasID := packet.Args["id"]
		role, hasRole := packet.Args["role"]

		if err != nil || packet.Command != initPacket || !hasID || !hasRole {
			socketErrorAndTerminate(conn, "could not decode init packet")
			return
		}

		// create new Client connection
		existingClient := s.locateClient(id)

		if existingClient != nil {
			// this should never happen, sanity check
			socketErrorAndTerminate(conn, "id already exists")
			return
		}

		client := Client{socket: conn, Role: role, Id: id}

		s.connections = append(
			s.connections,
			client,
		)

		// call OnConnect for all controllers
		for _, controller := range controllers {
			controller.OnConnect(&client)
		}

		// create a new Responder for this client
		responder := s.createResponder(client.Id)

		// enter read loop and block
		var readErr error
		for readErr == nil {
			var packet CommandPacket
			err := client.socket.ReadJSON(&packet)

			if err != nil {
				// read error, assume disconnect
				readErr = err
				// call OnDisconnect for all controllers
				s.disconnectClientById(client.Id, controllers)
			} else {
				// append client's id to FromId if not present
				if packet.FromId == "" {
					packet.FromId = client.Id
				}

				// parse command and add responder
				command := CommandFromPacket(packet)
				command.Responder = responder

				if command.TargetDeviceId != "" {
					// command is intended for another device, attempt to proxy it
					targetClient := s.locateClient(command.TargetDeviceId)
					if targetClient != nil {
						// target client found, proxy command
						targetClient.SendCommand(command)
					} else {
						responder.RespondError(
							"unable to transmit command to destination - no device with id '" + command.TargetDeviceId + "' exists",
						)
					}
				} else {
					// command is intended for server, pass to controller
					for _, controller := range controllers {
						controller.ParseCommand(&client, command)
					}
				}

			}

		}
	}
}

func socketErrorAndTerminate(conn *ws.Conn, message string) {
	_ = conn.WriteMessage(ws.TextMessage, []byte("error - '"+message+"', closing connection"))
	conn.Close()
}

func (s *WebSocketServer) disconnectClientById(id string, controllers []entities.Controller) {
	var newConnections []Client

	// remove this Client from list of connections
	for _, c := range s.connections {
		if c.Id != id {
			newConnections = append(newConnections, c)
		} else {
			c.socket.Close()

			for _, controller := range controllers {
				// call DisconnectEventHandler for all controllers
				controller.OnDisconnect(&c)
			}
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

func (s *WebSocketServer) createResponder(id string) entities.Responder {
	return entities.Responder{
		Respond: func(command entities.Command) {
			// fetch client
			client := s.locateClient(id)

			if client != nil {
				client.SendCommand(command)
			} else {
				panic("client doesn't exist")
			}
		},
		RespondError: func(error string) {
			client := s.locateClient(id)

			if client != nil {
				client.SendError(error)
			} else {
				panic("client doesn't exist")
			}
		},
	}
}
