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

package main

import (
	"os"
	"testing"
	"time"

	"github.com/livekit/signal-proxy/pkg/config"
	"github.com/livekit/signal-proxy/pkg/server"
	"github.com/livekit/signal-proxy/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	docker := test_utils.NewDocker("../test/docker/docker-compose.yml")
	err := docker.Up()
	if err != nil {
		println("Failed to start docker")
		os.Exit(1)
	}

	// Run the tests
	code := m.Run()

	os.Exit(code)
}

func Test_BasicNoProxy(t *testing.T) {
	participant, err := test_utils.NewLiveKitParticipant(7880, "../test/media/audio_track.ogg", false)
	defer participant.Disconnect()
	require.NoError(t, err, "participant create successfully")

	err = participant.ConnectAndPublish()
	assert.NoError(t, err, "participant should connect successfully")
}

func Test_HappyProxy(t *testing.T) {
	proxy := server.NewServer(&config.Config{DestinationLiveKitURL: "ws://localhost:7880", Port: 9001})
	defer proxy.Stop()

	participant, err := test_utils.NewLiveKitParticipant(9001, "../test/media/audio_track.ogg", false)
	defer participant.Disconnect()
	require.NoError(t, err, "participant create successfully")

	go func() {
		err := proxy.Run()
		assert.NoError(t, err, "proxy server should run successfully")
	}()

	time.Sleep(1 * time.Second)
	err = participant.ConnectAndPublish()
	assert.NoError(t, err, "participant should connect successfully")
}

func Test_ForceRelayHappy(t *testing.T) {
	proxy := server.NewServer(&config.Config{
		DestinationLiveKitURL: "ws://localhost:7880",
		Port:                  9002,
		ICEServers: []config.ICEServer{
			{
				Urls:       []string{"turn:127.0.0.1:3478?transport=udp"},
				Username:   "foo",
				Credential: "bar",
			},
			{
				Urls: []string{"stun:stun.l.google.com:19302"},
			},
		}})
	defer proxy.Stop()

	participant, err := test_utils.NewLiveKitParticipant(9002, "../test/media/audio_track.ogg", true)
	defer participant.Disconnect()
	require.NoError(t, err, "participant create successfully")

	go func() {
		err := proxy.Run()
		assert.NoError(t, err, "proxy server should run successfully")
	}()

	time.Sleep(1 * time.Second)
	err = participant.ConnectAndPublish()
	assert.NoError(t, err, "participant should connect successfully")
}

func Test_ForceRelayNoTurn(t *testing.T) {
	proxy := server.NewServer(&config.Config{DestinationLiveKitURL: "ws://localhost:7880", Port: 9003})
	defer proxy.Stop()

	participant, err := test_utils.NewLiveKitParticipant(9003, "../test/media/audio_track.ogg", true)
	defer participant.Disconnect()
	require.NoError(t, err, "participant create successfully")

	go func() {
		err := proxy.Run()
		assert.NoError(t, err, "proxy server should run successfully")
	}()

	time.Sleep(1 * time.Second)
	err = participant.ConnectAndPublish()
	assert.Error(t, err, "participant should not connect successfully")
}
