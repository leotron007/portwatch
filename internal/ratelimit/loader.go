package ratelimit

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds YAML-deserializable configuration for a rate limiter.
type Config struct {
	// Window is the rolling time window duration string (e.g. "1m", "30s").
	Window string `yaml:"window"`
	// Max is the maximum number of events allowed within the window.
	Max int `yaml:"max"`
}

// DefaultConfig returns a Config with sane defaults:
// a 1-minute window allowing up to 60 events.
func DefaultConfig() Config {
	return Config{
		Window: "1m",
		Max:    60,
	}
}

// FromConfig constructs a Limiter from a Config.
// Returns an error if the window cannot be parsed or the limiter is invalid.
func FromConfig(cfg Config) (*Limiter, error) {
	if cfg.Window == "" {
		cfg.Window = DefaultConfig().Window
	}
	if cfg.Max <= 0 {
		cfg.Max = DefaultConfig().Max
	}

	d, err := time.ParseDuration(cfg.Window)
	if err != nil {
		return nil, fmt.Errorf("ratelimit: invalid window %q: %w", cfg.Window, err)
	}

	return New(d, cfg.Max)
}

// FromBytes constructs a Limiter from YAML-encoded bytes.
func FromBytes(data []byte) (*Limiter, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ratelimit: unmarshal config: %w", err)
	}
	return FromConfig(cfg)
}
