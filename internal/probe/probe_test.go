package probe_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/probe"
)

func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	return port, func() { _ = ln.Close() }
}

func TestNew_DefaultHost(t *testing.T) {
	p, err := probe.New("", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}

func TestProbe_OpenPort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p, _ := probe.New("127.0.0.1", time.Second)
	r := p.Probe(port)

	if !r.Open {
		t.Fatalf("expected port %d to be open", port)
	}
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.Latency <= 0 {
		t.Fatal("expected positive latency")
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	p, _ := probe.New("127.0.0.1", 200*time.Millisecond)
	r := p.Probe(1)
	if r.Open {
		t.Fatal("expected port to be closed")
	}
	if r.Err == nil {
		t.Fatal("expected an error")
	}
}

func TestProbePorts_ReturnsAllResults(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p, _ := probe.New("127.0.0.1", time.Second)
	results := p.ProbePorts([]int{port, 1})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Open {
		t.Error("first port should be open")
	}
	if results[1].Open {
		t.Error("second port should be closed")
	}
}
