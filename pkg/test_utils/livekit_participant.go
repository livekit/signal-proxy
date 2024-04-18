// Copyright 2024 LiveKit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test_utils

import (
	"fmt"
	"time"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"github.com/pion/webrtc/v3"
)

type LiveKitParticipant struct {
	port       uint32
	audioTrack string
	room       *lksdk.Room
	forceRelay bool
}

func NewLiveKitParticipant(port uint32, audioTrack string, forceRelay bool) (*LiveKitParticipant, error) {
	return &LiveKitParticipant{
		port:       port,
		audioTrack: audioTrack,
		forceRelay: forceRelay,
	}, nil
}

func (p *LiveKitParticipant) ConnectAndPublish() error {
	url := fmt.Sprintf("ws://127.0.0.1:%d", p.port)
	apiKey := "devkey"
	apiSecret := "secret"
	roomName := "test-room"
	identity := "test-participant"
	roomCB := &lksdk.RoomCallback{
		OnDisconnected: func() {
			fmt.Println("Disconnected from room")
		},
	}
	fmt.Println("Connecting to room")
	iceTransportPolicy := webrtc.ICETransportPolicyAll
	if p.forceRelay {
		iceTransportPolicy = webrtc.ICETransportPolicyRelay
	}

	room, err := lksdk.ConnectToRoom(url, lksdk.ConnectInfo{
		APIKey:              apiKey,
		APISecret:           apiSecret,
		RoomName:            roomName,
		ParticipantIdentity: identity,
	}, roomCB, lksdk.WithICETransportPolicy(iceTransportPolicy))

	p.room = room

	if err != nil {
		return fmt.Errorf("failed to connect to room: %w", err)
	}

	fmt.Println("Connected to room")

	doneSignal := make(chan struct{})

	track, err := lksdk.NewLocalFileTrack(p.audioTrack,
		// control FPS to ensure synchronization
		lksdk.ReaderTrackWithFrameDuration(20*time.Millisecond),
		lksdk.ReaderTrackWithOnWriteComplete(func() {
			close(doneSignal)
		}),
	)

	if err != nil {
		return fmt.Errorf("failed to open track %w", err)
	}

	fmt.Print("Publishing track")
	if _, err = room.LocalParticipant.PublishTrack(track, &lksdk.TrackPublicationOptions{
		Name:   "audio_track",
		Source: livekit.TrackSource_MICROPHONE,
	}); err != nil {
		return err
	}
	fmt.Print("Published track")

	select {
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timed out waiting for track to finish")
	case <-doneSignal:
	}

	return nil
}

func (p *LiveKitParticipant) Disconnect() {
	if p.room != nil {
		p.room.Disconnect()
	}
}
