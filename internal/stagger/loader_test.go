package stagger

import (
	"testing"
	"time"
)

func TestDefaultConfig_HasSaneValues(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Count < 1 {
		t.Fatalf("default count = %d, want >= 1", cfg.Count)
	}
	if cfg.Window == "" {
		t.Fatal("default window is empty")
	}
}

func TestFromConfig_ValidConfig(t *testing.T) {
	cfg := Config{Count: 6, Window: "6s"}
	s, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count() != 6 {
		t.Fatalf("count = %d, want 6", s.Count())
	}
	if s.Window() != 6*time.Second {
		t.Fatalf("window = %v, want 6s", s.Window())
	}
}

func TestFromConfig_InvalidCount(t *testing.T) {
	_, err := FromConfig(Config{Count: 0, Window: "1s"})
	if err == nil {
		t.Fatal("expected error for count=0")
	}
}

func TestFromConfig_InvalidWindow(t *testing.T) {
	_, err := FromConfig(Config{Count: 2, Window: "not-a-duration"})
	if err == nil {
		t.Fatal("expected error for invalid window")
	}
}

func TestFromConfig_EmptyWindowUsesDefault(t *testing.T) {
	s, err := FromConfig(Config{Count: 2, Window: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Window() != time.Second {
		t.Fatalf("window = %v, want 1s", s.Window())
	}
}

func TestFromBytes_ValidYAML(t *testing.T) {
	data := []byte("count: 4\nwindow: 2s\n")
	s, err := FromBytes(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count() != 4 {
		t.Fatalf("count = %d, want 4", s.Count())
	}
}

func TestFromBytes_InvalidYAML(t *testing.T) {
	_, err := FromBytes([]byte(":::bad yaml:::"))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestFromBytes_MissingFieldsUseDefaults(t *testing.T) {
	s, err := FromBytes([]byte("{}"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count() < 1 {
		t.Fatalf("count = %d, want >= 1", s.Count())
	}
}
