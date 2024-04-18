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
