package budget

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the YAML-serialisable configuration for a Budget.
type Config struct {
	Capacity int    `yaml:"capacity"`
	Window   string `yaml:"window"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Capacity: 10,
		Window:   "1h",
	}
}

// FromConfig constructs a Budget from a Config.
func FromConfig(cfg Config) (*Budget, error) {
	if cfg.Capacity <= 0 {
		cfg.Capacity = DefaultConfig().Capacity
	}
	win := cfg.Window
	if win == "" {
		win = DefaultConfig().Window
	}
	d, err := time.ParseDuration(win)
	if err != nil {
		return nil, fmt.Errorf("budget: invalid window %q: %w", win, err)
	}
	if d <= 0 {
		return nil, errors.New("budget: window must be positive")
	}
	return New(cfg.Capacity, d)
}

// FromBytes parses YAML bytes and returns a Budget.
func FromBytes(data []byte) (*Budget, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("budget: yaml unmarshal: %w", err)
	}
	return FromConfig(cfg)
}
