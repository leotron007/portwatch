package throttle

import (
	"testing"
	"time"
)

func TestFromConfig_ValidCooldown(t *testing.T) {
	cfg := Config{CooldownSeconds: 30, MaxBurst: 2}
	th, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil Throttle")
	}
	if th.cooldown != 30*time.Second {
		t.Errorf("expected cooldown 30s, got %v", th.cooldown)
	}
	if th.maxBurst != 2 {
		t.Errorf("expected maxBurst 2, got %d", th.maxBurst)
	}
}

func TestFromConfig_ZeroCooldown(t *testing.T) {
	cfg := Config{CooldownSeconds: 0, MaxBurst: 1}
	th, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.cooldown != 0 {
		t.Errorf("expected zero cooldown, got %v", th.cooldown)
	}
}

func TestFromConfig_NegativeCooldown(t *testing.T) {
	cfg := Config{CooldownSeconds: -1}
	_, err := FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for negative cooldown_seconds")
	}
}

func TestFromConfig_DefaultMaxBurst(t *testing.T) {
	cfg := Config{CooldownSeconds: 5, MaxBurst: 0}
	th, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.maxBurst != 1 {
		t.Errorf("expected default maxBurst 1, got %d", th.maxBurst)
	}
}

func TestFromBytes_ValidYAML(t *testing.T) {
	yaml := []byte("cooldown_seconds: 10\nmax_burst: 3\n")
	th, err := FromBytes(yaml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.cooldown != 10*time.Second {
		t.Errorf("expected 10s cooldown, got %v", th.cooldown)
	}
	if th.maxBurst != 3 {
		t.Errorf("expected maxBurst 3, got %d", th.maxBurst)
	}
}

func TestFromBytes_InvalidYAML(t *testing.T) {
	_, err := FromBytes([]byte(":::not valid yaml:::"))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
