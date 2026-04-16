package probe

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds YAML-deserializable probe settings.
type Config struct {
	Host    string        `yaml:"host"`
	Timeout time.Duration `yaml:"timeout"`
}

// FromBytes parses YAML bytes into a Prober.
func FromBytes(data []byte) (*Prober, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("probe: parse config: %w", err)
	}
	return FromConfig(cfg)
}

// FromConfig builds a Prober from a Config struct.
func FromConfig(cfg Config) (*Prober, error) {
	return New(cfg.Host, cfg.Timeout)
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:    "127.0.0.1",
		Timeout: 2 * time.Second,
	}
}
