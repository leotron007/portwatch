// Package fingerprint provides port-fingerprinting utilities that
// record a lightweight signature (protocol hint + banner snippet) for
// each open port so that portwatch can detect service-level changes
// beyond simple open/closed transitions.
package fingerprint

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// DefaultTimeout is used when no explicit deadline is supplied.
const DefaultTimeout = 2 * time.Second

// Fingerprint holds the captured metadata for a single port.
type Fingerprint struct {
	Port    int    `json:"port"`
	Banner  string `json:"banner,omitempty"`
	Protocol string `json:"protocol"`
}

// Prober grabs a banner from a TCP port and returns a Fingerprint.
type Prober struct {
	Host    string
	Timeout time.Duration
}

// New returns a Prober for the given host.  A zero timeout falls back
// to DefaultTimeout.
func New(host string, timeout time.Duration) *Prober {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return &Prober{Host: host, Timeout: timeout}
}

// Probe connects to the given port, reads up to 256 bytes of banner
// data and returns a Fingerprint.  Errors during read are silently
// swallowed — a partial banner is still useful.
func (p *Prober) Probe(port int) (Fingerprint, error) {
	addr := fmt.Sprintf("%s:%d", p.Host, port)
	conn, err := net.DialTimeout("tcp", addr, p.Timeout)
	if err != nil {
		return Fingerprint{}, fmt.Errorf("fingerprint: dial %s: %w", addr, err)
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(p.Timeout))

	buf := make([]byte, 256)
	n, _ := conn.Read(buf)

	banner := strings.TrimSpace(string(buf[:n]))
	banner = sanitize(banner)

	return Fingerprint{
		Port:     port,
		Banner:   banner,
		Protocol: guessProtocol(banner),
	}, nil
}

// sanitize removes non-printable characters from a raw banner.
func sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 0x20 && r < 0x7F {
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}

// guessProtocol returns a coarse protocol hint based on banner content.
func guessProtocol(banner string) string {
	upper := strings.ToUpper(banner)
	switch {
	case strings.HasPrefix(upper, "SSH"):
		return "ssh"
	case strings.HasPrefix(upper, "HTTP"), strings.HasPrefix(upper, "GET "):
		return "http"
	case strings.HasPrefix(upper, "220 ") || strings.Contains(upper, "FTP"):
		return "ftp"
	case strings.HasPrefix(upper, "220") && strings.Contains(upper, "SMTP"):
		return "smtp"
	default:
		return "tcp"
	}
}
