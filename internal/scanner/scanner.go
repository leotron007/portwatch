package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port with its metadata.
type Port struct {
	Number   int
	Protocol string
	Address  string
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a new Scanner with the given host and timeout.
func New(host string, timeout time.Duration) *Scanner {
	if host == "" {
		host = "127.0.0.1"
	}
	return &Scanner{
		Host:    host,
		Timeout: timeout,
	}
}

// Scan checks the given port range and returns all open ports.
func (s *Scanner) Scan(startPort, endPort int) ([]Port, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	var open []Port
	for port := startPort; port <= endPort; port++ {
		if s.isOpen(port, "tcp") {
			open = append(open, Port{
				Number:   port,
				Protocol: "tcp",
				Address:  s.Host,
			})
		}
	}
	return open, nil
}

// ScanPorts checks a specific list of ports and returns all open ones.
func (s *Scanner) ScanPorts(ports []int) ([]Port, error) {
	var open []Port
	for _, port := range ports {
		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("invalid port number: %d", port)
		}
		if s.isOpen(port, "tcp") {
			open = append(open, Port{
				Number:   port,
				Protocol: "tcp",
				Address:  s.Host,
			})
		}
	}
	return open, nil
}

// isOpen attempts a TCP connection to determine if the port is open.
func (s *Scanner) isOpen(port int, proto string) bool {
	addr := fmt.Sprintf("%s:%d", s.Host, port)
	conn, err := net.DialTimeout(proto, addr, s.Timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
