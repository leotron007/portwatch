// Package config loads and validates the top-level portwatch configuration.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the complete portwatch runtime configuration.
type Config struct {
	// Scan defines port-scanning behaviour.
	Scan ScanConfig `yaml:"scan"`

	// Filter defines which ports are considered.
	Filter FilterConfig `yaml:"filter"`

	// Alert defines notification preferences.
	Alert AlertConfig `yaml:"alert"`
}

// ScanConfig controls how and how often ports are scanned.
type ScanConfig struct {
	Host     string        `yaml:"host"`
	PortMin  int           `yaml:"port_min"`
	PortMax  int           `yaml:"port_max"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

// FilterConfig mirrors internal/filter options.
type FilterConfig struct {
	IgnoredPorts []int `yaml:"ignored_ports"`
	AllowedPorts []int `yaml:"allowed_ports"`
}

// AlertConfig controls alert output.
type AlertConfig struct {
	Format string `yaml:"format"` // "text" | "json"
	Level  string `yaml:"level"`  // "info" | "warn" | "error"
}

// defaults applied when fields are zero-valued.
func applyDefaults(c *Config) {
	if c.Scan.Host == "" {
		c.Scan.Host = "127.0.0.1"
	}
	if c.Scan.PortMin == 0 {
		c.Scan.PortMin = 1
	}
	if c.Scan.PortMax == 0 {
		c.Scan.PortMax = 65535
	}
	if c.Scan.Interval == 0 {
		c.Scan.Interval = 30 * time.Second
	}
	if c.Scan.Timeout == 0 {
		c.Scan.Timeout = 500 * time.Millisecond
	}
	if c.Alert.Format == "" {
		c.Alert.Format = "text"
	}
	if c.Alert.Level == "" {
		c.Alert.Level = "info"
	}
}

// validate returns an error if the configuration is logically inconsistent.
func validate(c *Config) error {
	if c.Scan.PortMin < 1 || c.Scan.PortMin > 65535 {
		return fmt.Errorf("scan.port_min %d out of range [1, 65535]", c.Scan.PortMin)
	}
	if c.Scan.PortMax < 1 || c.Scan.PortMax > 65535 {
		return fmt.Errorf("scan.port_max %d out of range [1, 65535]", c.Scan.PortMax)
	}
	if c.Scan.PortMin > c.Scan.PortMax {
		return fmt.Errorf("scan.port_min (%d) must be <= scan.port_max (%d)", c.Scan.PortMin, c.Scan.PortMax)
	}
	switch c.Alert.Format {
	case "text", "json":
	default:
		return fmt.Errorf("alert.format %q must be \"text\" or \"json\"", c.Alert.Format)
	}
	return nil
}

// LoadFromFile reads a YAML file and returns a validated Config.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}
	return LoadFromBytes(data)
}

// LoadFromBytes parses YAML bytes and returns a validated Config.
func LoadFromBytes(data []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}
	applyDefaults(&c)
	if err := validate(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Default returns a Config populated entirely with default values.
func Default() *Config {
	var c Config
	applyDefaults(&c)
	return &c
}
