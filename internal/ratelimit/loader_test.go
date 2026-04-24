package ratelimit

import (
	"testing"
	"time"
)

func TestDefaultConfig_HasSaneValues(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Window != "1m" {
		t.Errorf("expected window 1m, got %s", cfg.Window)
	}
	if cfg.Max != 60 {
		t.Errorf("expected max 60, got %d", cfg.Max)
	}
}

func TestFromConfig_ValidConfig(t *testing.T) {
	cfg := Config{Window: "30s", Max: 10}
	l, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestFromConfig_InvalidWindow(t *testing.T) {
	cfg := Config{Window: "not-a-duration", Max: 5}
	_, err := FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid window")
	}
}

func TestFromConfig_ZeroMaxUsesDefault(t *testing.T) {
	cfg := Config{Window: "1m", Max: 0}
	l, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Exhaust default max (60) calls within window — all should pass.
	now := time.Now()
	for i := 0; i < 60; i++ {
		if !l.Allow(now) {
			t.Fatalf("call %d should be allowed", i+1)
		}
	}
}

func TestFromConfig_EmptyWindowUsesDefault(t *testing.T) {
	cfg := Config{Window: "", Max: 5}
	l, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestFromBytes_ValidYAML(t *testing.T) {
	data := []byte("window: 10s\nmax: 3\n")
	l, err := FromBytes(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestFromBytes_InvalidYAML(t *testing.T) {
	data := []byte(": : invalid yaml:::")
	_, err := FromBytes(data)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestFromBytes_InvalidWindow(t *testing.T) {
	data := []byte("window: bad\nmax: 5\n")
	_, err := FromBytes(data)
	if err == nil {
		t.Fatal("expected error for invalid window in YAML")
	}
}
