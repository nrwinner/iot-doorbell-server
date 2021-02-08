package commands

import (
	"doorbell-server/src/modules/websocket"
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v3"
)

const (
	offer_command     = "webrtc/offer"
	answer_command    = "webrtc/answer"
	candidate_command = "webrtc/candidate"
)

var connections = make(map[string]*webrtc.PeerConnection)

func ParseCommand(client websocket.Client, command string, args []string) {
	fmt.Println(command)
	if command == offer_command {
		var offer webrtc.SessionDescription
		err := json.Unmarshal([]byte(args[0]), &offer)
		if err != nil {
			panic(err)
		}

		peer := connections[client.Id]
		if peer == nil {
			peer = createPeerConnection(client.Id)
		}

		// set remote description to offer from client
		err = peer.SetRemoteDescription(offer)
		if err != nil {
			panic(err)
		}

		// create an answer
		answer, err := peer.CreateAnswer(nil)
		if err != nil {
			panic(err)
		}

		// set answer as local description
		err = peer.SetLocalDescription(answer)
		if err != nil {
			panic(err)
		}

		// send answer to client via signal server
		answerStr, err := json.Marshal(answer)
		if err != nil {
			panic(err)
		}

		packet := websocket.CommandPacket{
			Packet: websocket.Packet{
				Id:         client.Id,
				PacketType: "command",
			},
			Command: answer_command,
			Args:    []string{string(answerStr)},
		}

		client.SendCommand(packet)
	} else if command == candidate_command {
		peer := connections[client.Id]

		if peer == nil {
			peer = createPeerConnection(client.Id)
		}

		var candidate webrtc.ICECandidateInit
		err := json.Unmarshal([]byte(args[0]), &candidate)
		if err != nil {
			panic(err)
		}

		err = peer.AddICECandidate(candidate)
		if err != nil {
			panic(err)
		}
	}
}

func createPeerConnection(id string) *webrtc.PeerConnection {
	var err error
	connections[id], err = webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	connections[id].OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state.String() == "disconnected" {
			connections[id] = nil
		}
		fmt.Println(state)
	})

	return connections[id]
}
