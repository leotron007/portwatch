package config

import (
	"testing"
	"time"
)

func TestDefault_HasSaneValues(t *testing.T) {
	c := Default()
	if c.Scan.Host != "127.0.0.1" {
		t.Errorf("default host = %q, want 127.0.0.1", c.Scan.Host)
	}
	if c.Scan.PortMin != 1 {
		t.Errorf("default port_min = %d, want 1", c.Scan.PortMin)
	}
	if c.Scan.PortMax != 65535 {
		t.Errorf("default port_max = %d, want 65535", c.Scan.PortMax)
	}
	if c.Scan.Interval != 30*time.Second {
		t.Errorf("default interval = %v, want 30s", c.Scan.Interval)
	}
	if c.Alert.Format != "text" {
		t.Errorf("default format = %q, want text", c.Alert.Format)
	}
}

func TestLoadFromBytes_ValidYAML(t *testing.T) {
	yaml := []byte(`
scan:
  host: "0.0.0.0"
  port_min: 1024
  port_max: 9000
  interval: 10s
alert:
  format: json
  level: warn
`)
	c, err := LoadFromBytes(yaml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Scan.Host != "0.0.0.0" {
		t.Errorf("host = %q, want 0.0.0.0", c.Scan.Host)
	}
	if c.Scan.PortMin != 1024 {
		t.Errorf("port_min = %d, want 1024", c.Scan.PortMin)
	}
	if c.Alert.Format != "json" {
		t.Errorf("format = %q, want json", c.Alert.Format)
	}
	if c.Scan.Interval != 10*time.Second {
		t.Errorf("interval = %v, want 10s", c.Scan.Interval)
	}
}

func TestLoadFromBytes_InvalidYAML(t *testing.T) {
	_, err := LoadFromBytes([]byte(":::bad yaml:::"))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestValidate_InvalidPortRange(t *testing.T) {
	yaml := []byte(`
scan:
  port_min: 9000
  port_max: 1000
`)
	_, err := LoadFromBytes(yaml)
	if err == nil {
		t.Fatal("expected error for inverted port range, got nil")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	yaml := []byte(`
alert:
  format: xml
`)
	_, err := LoadFromBytes(yaml)
	if err == nil {
		t.Fatal("expected error for unsupported alert format, got nil")
	}
}

func TestValidate_PortOutOfRange(t *testing.T) {
	yaml := []byte(`
scan:
  port_min: 0
  port_max: 1000
`)
	_, err := LoadFromBytes(yaml)
	if err == nil {
		t.Fatal("expected error for port_min=0, got nil")
	}
}
