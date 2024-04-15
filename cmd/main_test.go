package main

import (
	"testing"

	"github.com/livekit/signal-proxy/pkg/config"
	"github.com/livekit/signal-proxy/pkg/server"
	"github.com/livekit/signal-proxy/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BasicNoProxy(t *testing.T) {
	docker := test_utils.NewDocker("../test/docker/docker-compose.yml")
	err := docker.Up()
	defer docker.Down()
	require.NoError(t, err, "docker up should succeed")

	participant, err := test_utils.NewLiveKitParticipant(7880, "../test/media/audio_track.ogg", false)
	defer participant.Disconnect()
	require.NoError(t, err, "participant create successfully")

	err = participant.ConnectAndPublish()
	assert.NoError(t, err, "participant should connect successfully")
}

func Test_HappyProxy(t *testing.T) {
	docker := test_utils.NewDocker("../test/docker/docker-compose.yml")
	err := docker.Up()
	defer docker.Down()
	require.NoError(t, err, "docker up should succeed")

	proxy := server.NewServer(&config.Config{DestinationHost: "localhost:7880", Port: 9000})
	defer proxy.Stop()

	participant, err := test_utils.NewLiveKitParticipant(9000, "../test/media/audio_track.ogg", false)
	defer participant.Disconnect()
	require.NoError(t, err, "participant create successfully")

	go func() {
		err := proxy.Run()
		assert.NoError(t, err, "proxy server should run successfully")
	}()

	err = participant.ConnectAndPublish()
	assert.NoError(t, err, "participant should connect successfully")
}

func Test_ForceRelayHappy(t *testing.T) {
	docker := test_utils.NewDocker("../test/docker/docker-compose.yml")
	err := docker.Up()
	defer docker.Down()
	require.NoError(t, err, "docker up should succeed")

	proxy := server.NewServer(&config.Config{DestinationHost: "localhost:7880", Port: 9000})
	defer proxy.Stop()

	participant, err := test_utils.NewLiveKitParticipant(9000, "../test/media/audio_track.ogg", true)
	defer participant.Disconnect()
	require.NoError(t, err, "participant create successfully")

	go func() {
		err := proxy.Run()
		assert.NoError(t, err, "proxy server should run successfully")
	}()

	err = participant.ConnectAndPublish()
	assert.NoError(t, err, "participant should connect successfully")
}

func Test_ForceRelayNoTurn(t *testing.T) {
	docker := test_utils.NewDocker("../test/docker/docker-compose.yml")
	err := docker.Up()
	defer docker.Down()
	require.NoError(t, err, "docker up should succeed")

	proxy := server.NewServer(&config.Config{DestinationHost: "localhost:7880", Port: 9000})
	defer proxy.Stop()

	participant, err := test_utils.NewLiveKitParticipant(9000, "../test/media/audio_track.ogg", true)
	defer participant.Disconnect()
	require.NoError(t, err, "participant create successfully")

	go func() {
		err := proxy.Run()
		assert.NoError(t, err, "proxy server should run successfully")
	}()

	err = participant.ConnectAndPublish()
	assert.Error(t, err, "participant should not connect successfully")
}
