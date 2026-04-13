package fingerprint_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/fingerprint"
)

// startBannerListener opens a TCP listener that writes banner on accept.
func startBannerListener(t *testing.T, banner string) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			fmt.Fprint(conn, banner)
			conn.Close()
		}
	}()

	return ln.Addr().(*net.TCPAddr).Port
}

func TestNew_DefaultHost(t *testing.T) {
	p := fingerprint.New("", 0)
	if p.Host != "127.0.0.1" {
		t.Errorf("expected default host 127.0.0.1, got %q", p.Host)
	}
	if p.Timeout != fingerprint.DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", fingerprint.DefaultTimeout, p.Timeout)
	}
}

func TestProbe_ReturnsBanner(t *testing.T) {
	port := startBannerListener(t, "SSH-2.0-OpenSSH_8.9")
	p := fingerprint.New("127.0.0.1", time.Second)

	fp, err := p.Probe(port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fp.Banner != "SSH-2.0-OpenSSH_8.9" {
		t.Errorf("banner mismatch: got %q", fp.Banner)
	}
	if fp.Protocol != "ssh" {
		t.Errorf("expected protocol ssh, got %q", fp.Protocol)
	}
	if fp.Port != port {
		t.Errorf("expected port %d, got %d", port, fp.Port)
	}
}

func TestProbe_HTTPProtocol(t *testing.T) {
	port := startBannerListener(t, "HTTP/1.1 200 OK")
	p := fingerprint.New("127.0.0.1", time.Second)

	fp, err := p.Probe(port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fp.Protocol != "http" {
		t.Errorf("expected protocol http, got %q", fp.Protocol)
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	p := fingerprint.New("127.0.0.1", 200*time.Millisecond)
	_, err := p.Probe(1) // port 1 is almost certainly closed
	if err == nil {
		t.Error("expected error probing closed port, got nil")
	}
}

func TestProbe_NoBannerDefaultsTCP(t *testing.T) {
	// Listener that closes immediately without writing.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	port := ln.Addr().(*net.TCPAddr).Port
	p := fingerprint.New("127.0.0.1", time.Second)
	fp, err := p.Probe(port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fp.Protocol != "tcp" {
		t.Errorf("expected fallback protocol tcp, got %q", fp.Protocol)
	}
}
