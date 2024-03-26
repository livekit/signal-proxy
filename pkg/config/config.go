package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DestinationHost string      `yaml:"destinationHost"`
	AllowedHosts    []string    `yaml:"allowedHosts"`
	Port            uint32      `yaml:"port"`
	ICEServers      []ICEServer `yaml:"iceServers"`
}

type ICEServer struct {
	Urls       []string `yaml:"urls,omitempty"`
	Username   string   `yaml:"username,omitempty"`
	Credential string   `yaml:"credential,omitempty"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("LK_CONFIG_PATH")
	if configPath == "" {
		configPath = "."
	}

	viper.SetConfigName("config")   // Name of config file (without extension)
	viper.SetConfigType("yaml")     // or viper.SetConfigType("YAML")
	viper.AddConfigPath(configPath) // Path to look for the config file in

	// Set default values
	viper.SetDefault("allowedHosts", []string{"0.0.0.0"})

	// Automatic binding of environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LK") // Prefix for environment variables

	var cfg Config

	err := viper.Unmarshal(&cfg)
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
	if c.DestinationHost == "" {
		return fmt.Errorf("destinationHost cannot be empty")
	}
	// Add custom validation logic here
	return nil
}
