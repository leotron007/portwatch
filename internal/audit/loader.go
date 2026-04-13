package audit

import (
	"fmt"
	"io"
	"os"
)

// Config holds configuration for the audit logger.
type Config struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// FromConfig returns a Logger and its underlying WriteCloser based on cfg.
// When cfg.Enabled is false or cfg.Path is empty, the logger writes to
// io.Discard and the returned closer is a no-op.
func FromConfig(cfg Config) (*Logger, io.Closer, error) {
	if !cfg.Enabled || cfg.Path == "" {
		return New(io.Discard), io.NopCloser(nil), nil
	}

	w, err := FileWriter(cfg.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("audit: from config: %w", err)
	}

	return New(w), w, nil
}

// DefaultConfig returns a Config with audit logging disabled.
func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Path:    "/var/log/portwatch/audit.log",
	}
}

// nopCloser wraps a nil closer so callers can always call Close.
type nopCloser struct{}

func (nopCloser) Close() error { return nil }

func init() {
	// Ensure os package is used (FileWriter depends on it).
	_ = os.DevNull
}
