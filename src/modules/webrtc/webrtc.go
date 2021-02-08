package webrtc

import (
	"github.com/pion/webrtc/v3"
)

func CreatePeerConnection(offer webrtc.SessionDescription) (*webrtc.PeerConnection, webrtc.SessionDescription) {
	peer, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	err = peer.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	answer, err := peer.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	err = peer.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	return peer, answer
}
