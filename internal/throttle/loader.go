package throttle

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the YAML-decoded throttle configuration.
type Config struct {
	// CooldownSeconds is the minimum number of seconds between repeated
	// alerts for the same port/event key. Zero disables throttling.
	CooldownSeconds int `yaml:"cooldown_seconds"`

	// MaxBurst is the maximum number of events allowed before the cooldown
	// kicks in. Defaults to 1 when not set.
	MaxBurst int `yaml:"max_burst"`
}

// FromConfig constructs a Throttle from a Config struct.
func FromConfig(cfg Config) (*Throttle, error) {
	if cfg.CooldownSeconds < 0 {
		return nil, fmt.Errorf("throttle: cooldown_seconds must be >= 0, got %d", cfg.CooldownSeconds)
	}
	maxBurst := cfg.MaxBurst
	if maxBurst <= 0 {
		maxBurst = 1
	}
	cooldown := time.Duration(cfg.CooldownSeconds) * time.Second
	return New(cooldown, maxBurst), nil
}

// FromBytes parses YAML bytes into a Throttle.
func FromBytes(data []byte) (*Throttle, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("throttle: failed to parse YAML: %w", err)
	}
	return FromConfig(cfg)
}
