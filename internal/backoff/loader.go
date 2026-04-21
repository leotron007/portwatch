package backoff

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the YAML-serialisable configuration for a Backoff.
type Config struct {
	Base   time.Duration `yaml:"base"`
	Max    time.Duration `yaml:"max"`
	Factor float64       `yaml:"factor"`
	Jitter float64       `yaml:"jitter"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Base:   500 * time.Millisecond,
		Max:    30 * time.Second,
		Factor: 2.0,
		Jitter: 0.1,
	}
}

// FromConfig constructs a Backoff from a Config, filling zero values with
// defaults before validation.
func FromConfig(cfg Config) (*Backoff, error) {
	def := DefaultConfig()
	if cfg.Base == 0 {
		cfg.Base = def.Base
	}
	if cfg.Max == 0 {
		cfg.Max = def.Max
	}
	if cfg.Factor == 0 {
		cfg.Factor = def.Factor
	}
	return New(cfg.Base, cfg.Max, cfg.Factor, cfg.Jitter)
}

// FromBytes parses YAML bytes into a Config and delegates to FromConfig.
func FromBytes(data []byte) (*Backoff, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("backoff: parse config: %w", err)
	}
	return FromConfig(cfg)
}
