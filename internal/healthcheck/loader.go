package healthcheck

import (
	"fmt"
	"os"
	"time"
)

// Config holds optional healthcheck configuration.
type Config struct {
	// Interval between automatic health writes (0 disables automatic writes).
	IntervalSeconds int `yaml:"interval_seconds"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{IntervalSeconds: 60}
}

// FromConfig validates cfg and returns a configured Reporter writing to w.
// If w is nil, os.Stdout is used.
func FromConfig(cfg Config, w interface{ Write([]byte) (int, error) }) (*Reporter, time.Duration, error) {
	if cfg.IntervalSeconds < 0 {
		return nil, 0, fmt.Errorf("healthcheck: interval_seconds must be >= 0, got %d", cfg.IntervalSeconds)
	}
	var interval time.Duration
	if cfg.IntervalSeconds > 0 {
		interval = time.Duration(cfg.IntervalSeconds) * time.Second
	}
	var writer interface{ Write([]byte) (int, error) } = os.Stdout
	if w != nil {
		writer = w
	}
	return New(writer), interval, nil
}
