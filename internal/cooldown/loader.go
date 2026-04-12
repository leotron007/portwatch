package cooldown

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds YAML-deserializable cooldown configuration.
type Config struct {
	// CooldownSeconds is the per-port quiet period in seconds.
	CooldownSeconds int `yaml:"cooldown_seconds"`
}

// FromConfig constructs a Tracker from a Config.
// Negative values are treated as zero (no cooldown).
func FromConfig(cfg Config) (*Tracker, error) {
	if cfg.CooldownSeconds < 0 {
		cfg.CooldownSeconds = 0
	}
	return New(time.Duration(cfg.CooldownSeconds) * time.Second), nil
}

// FromBytes parses YAML bytes and returns a configured Tracker.
func FromBytes(data []byte) (*Tracker, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cooldown: parse config: %w", err)
	}
	return FromConfig(cfg)
}
