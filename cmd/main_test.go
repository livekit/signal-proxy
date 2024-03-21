package main

import (
	"sync"
	"testing"

	"github.com/livekit/signal-proxy/pkg/config"
	"github.com/livekit/signal-proxy/pkg/server"
	"github.com/livekit/signal-proxy/pkg/test_utils"
)

func TestE2E(t *testing.T) {
	// Create two livekit servers, a proxy server, and a participant
	server1, err := test_utils.NewLiveKitServer(7001)
	if err != nil {
		t.Fatal(err)
	}

	// err = server1.Run()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	server2, err := test_utils.NewLiveKitServer(8001)
	if err != nil {
		t.Fatal(err)
	}

	proxy, err := server.NewServer(&config.Config{DestinationHost: "localhost:7001", Port: 9001})
	if err != nil {
		t.Fatal(err)
	}

	participant, err := test_utils.NewLiveKitParticipant(9001)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}

	errChan := make(chan error)
	done := make(chan struct{})

	// Start first livekit server
	wg.Add(1)
	go func() {
		defer wg.Done()
		s1Err := server1.Run()
		if err != nil {
			errChan <- s1Err
		}
	}()

	// Start second livekit server
	wg.Add(1)
	go func() {
		defer wg.Done()
		s2Err := server2.Run()
		if err != nil {
			errChan <- s2Err
		}
	}()

	// Start proxy server
	wg.Add(1)
	go func() {
		defer wg.Done()
		proxy.Run()
	}()

	// Start participant
	wg.Add(1)
	go func() {
		defer wg.Done()
		pErr := participant.RunAudioPublisher()
		if err != nil {
			errChan <- pErr
		}
	}()

	// Wait for everything
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errChan:
		t.Fatal(err)
	case <-done:
		// Everything has finished
	}
}
