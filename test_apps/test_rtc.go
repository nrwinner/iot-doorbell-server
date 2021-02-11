package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

const (
	ID = "doorbell01"
)

var connectionId string

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
		initPacket["Id"] = ID

		if err != nil {
			panic(err)
		}

		// send init packet
		conn.WriteJSON(initPacket)

		signal(conn, "new", nil)

		var newConnectionResult CommandPacket
		err := conn.ReadJSON(&newConnectionResult)
		connectionId = newConnectionResult.Args[0]

		peer, err := webrtc.NewPeerConnection(webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
		})
		if err != nil {
			panic(err)
		}

		peer.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
			fmt.Println("STATE", state)
		})

		peer.OnICECandidate(func(candidate *webrtc.ICECandidate) {
			// transmit the new ice candidate via the signal servers
			if candidate != nil {
				signal(conn, "candidate", candidate.ToJSON())
			}
		})

		peer.OnNegotiationNeeded(func() {
			go func() {
				// listen for answer packet and write into answerPacket variable
				var answerPacket CommandPacket
				done := false
				for !done {
					err = conn.ReadJSON(&answerPacket)
					if err != nil {
						panic(err)
					}

					if answerPacket.Command == "webrtc/answer" {
						done = true
					}
				}

				var answer webrtc.SessionDescription
				json.Unmarshal([]byte(answerPacket.Args[0]), &answer)

				err = peer.SetRemoteDescription(answer)
				if err != nil {
					panic(err)
				}
			}()

			// create offer
			offer, err := peer.CreateOffer(nil)
			if err != nil {
				panic(err)
			}

			err = peer.SetLocalDescription(offer)
			if err != nil {
				panic(err)
			}

			// send offer via signal server
			signal(conn, "offer", offer)
		})

		track, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "", "")
		peer.AddTrack(track)
	}()

	// block forever
	select {}
}

func signal(conn *websocket.Conn, name string, payload interface{}) {
	var args []string

	if connectionId != "" {
		args = append(args, connectionId)
	}

	if payload != nil {
		j, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}
		args = append(args, string(j))
	}

	packet := CommandPacket{
		Id:         ID,
		PacketType: "command",
		Command:    "webrtc/" + name,
		Args:       args,
	}

	err := conn.WriteJSON(packet)
	if err != nil {
		fmt.Println("Signal error for", name)
		panic(err)
	}
}
