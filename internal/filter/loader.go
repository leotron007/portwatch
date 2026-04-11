package filter

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Config is the YAML-serialisable representation of a Filter.
type Config struct {
	AllowedPorts []int `yaml:"allowed_ports"`
	IgnoredPorts []int `yaml:"ignored_ports"`
	MinPort      int   `yaml:"min_port"`
	MaxPort      int   `yaml:"max_port"`
}

// FromConfig converts a Config into a Filter, applying defaults for zero values.
func FromConfig(c Config) (*Filter, error) {
	f := New()
	if c.MinPort != 0 {
		f.MinPort = c.MinPort
	}
	if c.MaxPort != 0 {
		f.MaxPort = c.MaxPort
	}
	if f.MinPort > f.MaxPort {
		return nil, fmt.Errorf("filter: min_port %d exceeds max_port %d", f.MinPort, f.MaxPort)
	}
	f.AllowedPorts = c.AllowedPorts
	f.IgnoredPorts = c.IgnoredPorts
	return f, nil
}

// FromBytes parses YAML bytes into a Filter.
func FromBytes(data []byte) (*Filter, error) {
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("filter: yaml parse error: %w", err)
	}
	return FromConfig(c)
}
