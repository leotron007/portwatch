package filter

import (
	"testing"
)

func TestFromBytes_ValidYAML(t *testing.T) {
	yaml := []byte(`
min_port: 1024
max_port: 9000
ignored_ports: [22, 80]
allowed_ports: []
`)
	f, err := FromBytes(yaml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.MinPort != 1024 {
		t.Errorf("expected MinPort 1024, got %d", f.MinPort)
	}
	if f.MaxPort != 9000 {
		t.Errorf("expected MaxPort 9000, got %d", f.MaxPort)
	}
	if len(f.IgnoredPorts) != 2 {
		t.Errorf("expected 2 ignored ports, got %d", len(f.IgnoredPorts))
	}
}

func TestFromBytes_InvalidYAML(t *testing.T) {
	_, err := FromBytes([]byte(":::bad yaml:::"))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestFromConfig_InvalidRange(t *testing.T) {
	_, err := FromConfig(Config{MinPort: 9000, MaxPort: 1024})
	if err == nil {
		t.Error("expected error when min_port > max_port")
	}
}

func TestFromConfig_Defaults(t *testing.T) {
	f, err := FromConfig(Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.MinPort != 1 || f.MaxPort != 65535 {
		t.Errorf("expected defaults 1-65535, got %d-%d", f.MinPort, f.MaxPort)
	}
}

func TestFromConfig_AllowedPorts(t *testing.T) {
	f, err := FromConfig(Config{AllowedPorts: []int{443, 8443}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(443) {
		t.Error("expected 443 to be allowed")
	}
	if f.Allow(80) {
		t.Error("expected 80 to be excluded by allow-list")
	}
}

func TestFromConfig_IgnoredPorts(t *testing.T) {
	f, err := FromConfig(Config{IgnoredPorts: []int{22, 80, 443}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Allow(22) {
		t.Error("expected 22 to be ignored")
	}
	if f.Allow(80) {
		t.Error("expected 80 to be ignored")
	}
	if !f.Allow(8080) {
		t.Error("expected 8080 to be allowed (not in ignored list)")
	}
}
