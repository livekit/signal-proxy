package test_utils

import (
	"fmt"
	"os/exec"
)

type LiveKitParticipant struct {
	port uint32
	cmd  *exec.Cmd
}

func NewLiveKitParticipant(port uint32) (*LiveKitParticipant, error) {
	return &LiveKitParticipant{
		port: port,
	}, nil
}

func (p *LiveKitParticipant) RunAudioPublisher() error {
	p.cmd = exec.Command("livekit-cli",
		"load-test",
		"--audio-publishers",
		"1",
		"--url",
		fmt.Sprintf("http://localhost:%d", p.port),
		"--api-key",
		"devkey",
		"--api-secret",
		"secret",
	)

	stdoutPipe, err := p.cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %v\n", err)
		return err
	}

	stderrPipe, err := p.cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
		return err
	}

	go streamOutput(stdoutPipe, "livekit-cli-stdout")
	go streamOutput(stderrPipe, "livekit-cli-stdout")

	err = p.cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start livekit-cli: %w", err)
	}

	err = p.cmd.Wait()
	if err != nil {
		return fmt.Errorf("livekit-cli exited with error: %w", err)
	}

	return nil
}

func (p *LiveKitParticipant) Stop() {
	p.cmd.Process.Kill()
}
