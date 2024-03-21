package server

import (
	"fmt"
	"net/http"

	"github.com/livekit/signal-proxy/pkg/config"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) (*Server, error) {
	return &Server{
		cfg: cfg,
	}, nil
}

func (s *Server) Run() error {
	http.HandleFunc("/", s.handleConnection)
	http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), nil)
	return nil
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	connection, err := NewConnection(w, r, &s.cfg.DestinationHost)
	if err != nil {
		http.Error(w, "failed to create connection", http.StatusInternalServerError)
		return
	}

	connection.Run()
}
