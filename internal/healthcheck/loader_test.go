package healthcheck_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

func TestFromConfig_Defaults(t *testing.T) {
	cfg := healthcheck.DefaultConfig()
	if cfg.IntervalSeconds != 60 {
		t.Errorf("expected 60, got %d", cfg.IntervalSeconds)
	}
}

func TestFromConfig_ValidInterval(t *testing.T) {
	cfg := healthcheck.Config{IntervalSeconds: 30}
	r, interval, err := healthcheck.FromConfig(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
	if interval != 30*time.Second {
		t.Errorf("expected 30s, got %s", interval)
	}
}

func TestFromConfig_ZeroIntervalDisabled(t *testing.T) {
	cfg := healthcheck.Config{IntervalSeconds: 0}
	_, interval, err := healthcheck.FromConfig(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if interval != 0 {
		t.Errorf("expected 0 duration, got %s", interval)
	}
}

func TestFromConfig_NegativeIntervalError(t *testing.T) {
	cfg := healthcheck.Config{IntervalSeconds: -1}
	_, _, err := healthcheck.FromConfig(cfg, nil)
	if err == nil {
		t.Fatal("expected error for negative interval")
	}
}
