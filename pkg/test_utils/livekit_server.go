package test_utils

import (
	"fmt"
	"os/exec"

	"github.com/livekit/livekit-server/pkg/config"
	"github.com/livekit/mediatransportutil/pkg/rtcconfig"
	"gopkg.in/yaml.v3"
)

type LiveKitServer struct {
	port uint32
	cmd  *exec.Cmd
}

func NewLiveKitServer(port uint32) (*LiveKitServer, error) {
	return &LiveKitServer{
		port: port,
	}, nil
}

func (s *LiveKitServer) Run() error {
	rtcConf := rtcconfig.RTCConfig{
		UDPPort: rtcconfig.PortRange{
			Start: int(s.port + 1),
			End:   int(s.port + 2),
		},
		TCPPort: s.port + 3,
	}
	config := config.Config{
		Port:        s.port,
		Development: true,
		RTC:         config.RTCConfig{RTCConfig: rtcConf},
		Audio:       config.AudioConfig{},
		Video:       config.VideoConfig{},
		Room:        config.RoomConfig{},
		TURN:        config.TURNConfig{},
		LogLevel:    "debug",
	}
	configBytes, err := yaml.Marshal(&config)
	fmt.Print(string(configBytes))
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	s.cmd = exec.Command("livekit-server", "--dev", "--config-body", string(configBytes))

	stdoutPipe, err := s.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := s.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	go streamOutput(stdoutPipe, fmt.Sprintf("livekit-server-%d-stdout", s.port))
	go streamOutput(stderrPipe, fmt.Sprintf("livekit-server-%d-stderr", s.port))

	err = s.cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start livekit-server: %w", err)
	}
	err = s.cmd.Wait()
	if err != nil {
		return fmt.Errorf("livekit-server exited with error: %w", err)
	}
	return nil
}

func (s *LiveKitServer) Stop() {
	s.cmd.Process.Kill()
}
