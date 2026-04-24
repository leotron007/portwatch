package stagger

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the YAML-serialisable configuration for a Stagger.
type Config struct {
	Count  int    `yaml:"count"`
	Window string `yaml:"window"`
}

// DefaultConfig returns a safe default configuration.
func DefaultConfig() Config {
	return Config{
		Count:  1,
		Window: "1s",
	}
}

// FromConfig constructs a Stagger from a Config, using the real clock.
func FromConfig(cfg Config) (*Stagger, error) {
	if cfg.Count < 1 {
		return nil, errors.New("stagger: count must be >= 1")
	}
	raw := cfg.Window
	if raw == "" {
		raw = "1s"
	}
	w, err := time.ParseDuration(raw)
	if err != nil {
		return nil, fmt.Errorf("stagger: invalid window %q: %w", raw, err)
	}
	return New(cfg.Count, w, nil)
}

// FromBytes parses YAML bytes into a Stagger.
func FromBytes(data []byte) (*Stagger, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("stagger: unmarshal: %w", err)
	}
	if cfg.Count == 0 {
		cfg.Count = DefaultConfig().Count
	}
	if cfg.Window == "" {
		cfg.Window = DefaultConfig().Window
	}
	return FromConfig(cfg)
}
