package webrtc

import (
	"doorbell-server/src/entities"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
)

// TODO:NickW handle case where two entities want to directly connect, instead of using client/server model

func HandleNewCommand(controller *WebRTCController, command entities.Command) {
	// create a new uuid for the connection
	id := uuid.New().String()

	// sanity check: verify that the uuid doesn't already exist in map, regenerate until does not exist
	for {
		// if id doesn't exist in map, break out of loop
		if _, ok := controller.connections[id]; !ok {
			break
		}

		// replace id with new id and check again since it already exists
		id = uuid.New().String()
	}

	// create a new PeerConnection and store it in map, keyed by uuid
	var err error
	controller.connections[id], err = webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})

	if err != nil {
		panic(err)
	}

	// set disconnect handler for PeerConnection
	controller.connections[id].OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		fmt.Println("Initial Count", len(controller.connections))
		if state.String() == "disconnected" {
			// remove the id from the map of connections
			delete(controller.connections, id)
		}
		fmt.Println(state)
		fmt.Println("Final Count", len(controller.connections))
	})

	// respond to requester with id
	command.Responder.Respond(entities.Command{
		Path: NEW_CONFIRM_COMMAND,
		Args: []string{id},
	})
}

func HandleOfferCommand(controller *WebRTCController, command entities.Command) {
	// retrieve the peer id and offer from the command's Args array
	id := command.Args[0]
	offerStr := command.Args[1]

	// retrieve the existing peer from the controller, err if does not exist
	peer := controller.connections[id]

	if peer == nil {
		panic("no peer with id " + id)
	}

	// unmarshal the offer string into a SessionDescription instance
	var offer webrtc.SessionDescription
	err := json.Unmarshal([]byte(offerStr), &offer)
	if err != nil {
		panic(err)
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

	// respond to client with answer-command
	answerStr, err := json.Marshal(answer)
	if err != nil {
		panic(err)
	}

	answerCommand := entities.Command{
		Path: ANSWER_COMMAND,
		Args: []string{string(answerStr)},
	}

	command.Responder.Respond(answerCommand)
}

func HandleCandidateCommand(controller *WebRTCController, command entities.Command) {
	// fetch peer id and candidate from command's Args
	id := command.Args[0]
	candidateStr := command.Args[1]
	peer := controller.connections[id]

	if peer == nil {
		panic("no peer with id " + id)
	}

	// unmarshall the candidate string into an instance of ICECandidateInit
	var candidate webrtc.ICECandidateInit
	err := json.Unmarshal([]byte(candidateStr), &candidate)
	if err != nil {
		panic(err)
	}

	// add the decoded ICE candidate to the peer
	err = peer.AddICECandidate(candidate)
	if err != nil {
		panic(err)
	}
}
