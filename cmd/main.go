package main

import (
	"github.com/livekit/signal-proxy/pkg/config"
	"github.com/livekit/signal-proxy/pkg/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	server, err := server.NewServer(cfg)
	if err != nil {
		panic(err)
	}
	server.Run()
}
