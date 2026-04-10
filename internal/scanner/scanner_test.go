package scanner

import (
	"net"
	"testing"
	"time"
)

// startTestListener opens a TCP listener on an OS-assigned port and returns it.
func startTestListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestNew_DefaultHost(t *testing.T) {
	s := New("", time.Second)
	if s.Host != "127.0.0.1" {
		t.Errorf("expected default host 127.0.0.1, got %s", s.Host)
	}
}

func TestScan_DetectsOpenPort(t *testing.T) {
	ln, port := startTestListener(t)
	defer ln.Close()

	s := New("127.0.0.1", 500*time.Millisecond)
	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 || ports[0].Number != port {
		t.Errorf("expected port %d to be detected as open, got %v", port, ports)
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := New("127.0.0.1", 500*time.Millisecond)
	_, err := s.Scan(500, 100)
	if err == nil {
		t.Error("expected error for invalid range, got nil")
	}
}

func TestScan_ClosedPort(t *testing.T) {
	ln, port := startTestListener(t)
	ln.Close() // close immediately so the port is no longer open

	s := New("127.0.0.1", 200*time.Millisecond)
	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected no open ports, got %v", ports)
	}
}
