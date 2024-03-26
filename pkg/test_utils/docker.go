package test_utils

import (
	"fmt"
	"os/exec"
)

type Docker struct {
	cmd         *exec.Cmd
	composeFile string
}

func NewDocker(composeFile string) *Docker {
	return &Docker{composeFile: composeFile}
}

func (s *Docker) Up() error {
	s.cmd = exec.Command("docker-compose", "-f", s.composeFile, "up", "-d")
	// Pipe output
	stdoutPipe, err := s.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := s.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	go func() {
		streamOutput(stdoutPipe, "docker-compose-stdout")
	}()
	go func() {
		streamOutput(stderrPipe, "docker-compose-stderr")
	}()

	err = s.cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start docker compose: %w", err)
	}
	err = s.cmd.Wait()
	if err != nil {
		return fmt.Errorf("docker-compose exited with error: %w", err)
	}
	return nil
}

func (s *Docker) Down() {
	s.cmd = exec.Command("docker-compose", "-f", s.composeFile, "down")
	s.cmd.Run()
}
