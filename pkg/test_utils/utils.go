package test_utils

import (
	"bufio"
	"fmt"
	"io"
)

func streamOutput(pipe io.ReadCloser, label string) {
	reader := bufio.NewReader(pipe)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading %s: %v\n", label, err)
			}
			break
		}
		fmt.Printf("[%s] %s", label, line)
	}
}
