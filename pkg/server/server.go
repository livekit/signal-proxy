package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/livekit/signal-proxy/pkg/config"
)

type Server struct {
	cfg    *config.Config
	server *http.Server
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	s.server = &http.Server{Addr: fmt.Sprintf(":%d", s.cfg.Port), Handler: http.HandlerFunc(s.handleConnection)}
	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return fmt.Errorf("proxy server failed: %w", err)
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	connection, err := NewConnection(s.cfg, w, r)
	if err != nil {
		http.Error(w, "failed to create connection", http.StatusInternalServerError)
		return
	}

	connection.Run()
}
