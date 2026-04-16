// Package probe performs lightweight TCP connectivity checks and
// measures round-trip latency for a given host:port pair.
package probe

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	Port    int
	Open    bool
	Latency time.Duration
	Err     error
}

// Prober checks whether ports are reachable within a timeout.
type Prober struct {
	host    string
	timeout time.Duration
}

// New returns a Prober targeting host with the given dial timeout.
// If timeout is zero it defaults to 2 seconds.
func New(host string, timeout time.Duration) (*Prober, error) {
	if host == "" {
		host = "127.0.0.1"
	}
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Prober{host: host, timeout: timeout}, nil
}

// Probe dials the given port and returns a Result.
func (p *Prober) Probe(port int) Result {
	addr := fmt.Sprintf("%s:%d", p.host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	latency := time.Since(start)
	if err != nil {
		return Result{Port: port, Open: false, Latency: latency, Err: err}
	}
	_ = conn.Close()
	return Result{Port: port, Open: true, Latency: latency}
}

// ProbePorts probes each port in the slice and returns all results.
func (p *Prober) ProbePorts(ports []int) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		results = append(results, p.Probe(port))
	}
	return results
}
