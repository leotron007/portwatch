package graceperiod

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds YAML-serialisable settings for a Tracker.
type Config struct {
	// Window is the minimum duration a port must be observed before alerts
	// are emitted. Zero disables the grace period. Accepts Go duration strings
	// such as "5s" or "500ms".
	Window string `yaml:"window"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{Window: "3s"}
}

// FromConfig builds a Tracker from a Config.
func FromConfig(cfg Config) (*Tracker, error) {
	if cfg.Window == "" {
		cfg = DefaultConfig()
	}
	d, err := time.ParseDuration(cfg.Window)
	if err != nil {
		return nil, fmt.Errorf("graceperiod: invalid window %q: %w", cfg.Window, err)
	}
	return New(d, nil)
}

// FromBytes parses YAML and returns a Tracker.
func FromBytes(data []byte) (*Tracker, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("graceperiod: yaml: %w", err)
	}
	return FromConfig(cfg)
}
