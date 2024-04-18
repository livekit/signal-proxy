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

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DestinationLiveKitURL string      `yaml:"destination_livekit_url"`
	Port                  uint32      `yaml:"port"`
	ICEServers            []ICEServer `yaml:"ice_servers"`
}

type ICEServer struct {
	Urls       []string `yaml:"urls"`
	Username   string   `yaml:"username"`
	Credential string   `yaml:"credential"`
}

func LoadConfig() (*Config, error) {
	configFile := os.Getenv("LK_CONFIG_FILE")
	fmt.Println("neil configFile", configFile)
	if configFile == "" {
		return nil, fmt.Errorf("LK_CONFIG_FILE environment variable not set")
	}

	var cfg Config

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)

	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.DestinationLiveKitURL == "" {
		return fmt.Errorf("destination_livekit_url cannot be empty")
	}
	// Add custom validation logic here
	return nil
}
