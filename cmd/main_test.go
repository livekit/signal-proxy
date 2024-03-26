package main

import (
	"testing"

	"github.com/livekit/signal-proxy/pkg/config"
	"github.com/livekit/signal-proxy/pkg/server"
	"github.com/livekit/signal-proxy/pkg/test_utils"
	"github.com/stretchr/testify/assert"
)

func Test_HappyProxy(t *testing.T) {
	docker := test_utils.NewDocker("../test/docker/docker-compose-basic.yml")
	err := docker.Up()
	defer docker.Down()
	assert.NoError(t, err, "docker up should succeed")

	proxy := server.NewServer(&config.Config{DestinationHost: "localhost:7880", Port: 9000})
	defer proxy.Stop()

	participant, err := test_utils.NewLiveKitParticipant(9000, "../test/media/audio_track.ogg")
	defer participant.Disconnect()
	assert.NoError(t, err, "participant create successfully")

	go func() {
		err := proxy.Run()
		assert.NoError(t, err, "proxy server should run successfully")
	}()

	err = participant.ConnectAndPublish()
	assert.NoError(t, err, "participant should connect successfully")
}

// TODO
// func Test_ForceRelay(t *testing.T) {
// }
