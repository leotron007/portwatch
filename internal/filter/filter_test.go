package filter

import (
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	f := New()
	if f.MinPort != 1 || f.MaxPort != 65535 {
		t.Fatalf("expected default range 1-65535, got %d-%d", f.MinPort, f.MaxPort)
	}
}

func TestAllow_BelowMin(t *testing.T) {
	f := New()
	f.MinPort = 1024
	if f.Allow(80) {
		t.Error("expected port 80 to be rejected below MinPort")
	}
}

func TestAllow_AboveMax(t *testing.T) {
	f := New()
	f.MaxPort = 1024
	if f.Allow(8080) {
		t.Error("expected port 8080 to be rejected above MaxPort")
	}
}

func TestAllow_Ignored(t *testing.T) {
	f := New()
	f.IgnoredPorts = []int{22, 80}
	if f.Allow(22) {
		t.Error("expected port 22 to be ignored")
	}
	if f.Allow(80) {
		t.Error("expected port 80 to be ignored")
	}
	if !f.Allow(443) {
		t.Error("expected port 443 to be allowed")
	}
}

func TestAllow_AllowList(t *testing.T) {
	f := New()
	f.AllowedPorts = []int{443, 8080}
	if f.Allow(80) {
		t.Error("expected port 80 to be excluded by allow-list")
	}
	if !f.Allow(443) {
		t.Error("expected port 443 to be in allow-list")
	}
}

func TestApply_FiltersSlice(t *testing.T) {
	f := New()
	f.IgnoredPorts = []int{22}
	f.MaxPort = 9000
	input := []int{22, 80, 443, 9001}
	got := f.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d: %v", len(got), got)
	}
	if got[0] != 80 || got[1] != 443 {
		t.Errorf("unexpected ports: %v", got)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	f := New()
	got := f.Apply(nil)
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
