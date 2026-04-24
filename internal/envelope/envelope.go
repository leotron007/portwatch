// Package envelope wraps a scan event with routing metadata so that
// downstream consumers can filter, route, or log events without needing
// to inspect raw port numbers directly.
package envelope

import (
	"fmt"
	"time"
)

// Severity indicates how significant an event is.
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityCritical
)

// String returns a human-readable severity label.
func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "warning"
	case SeverityCritical:
		return "critical"
	default:
		return "info"
	}
}

// Event is the kind of port change that occurred.
type Event string

const (
	EventOpened Event = "opened"
	EventClosed Event = "closed"
)

// Envelope carries a port-change event together with routing metadata.
type Envelope struct {
	ID        string
	Host      string
	Port      int
	Event     Event
	Severity  Severity
	Timestamp time.Time
	Labels    map[string]string
}

// New creates an Envelope for the given host, port, and event type.
// The timestamp is set to the current UTC time.
func New(host string, port int, event Event, sev Severity) *Envelope {
	return &Envelope{
		ID:        fmt.Sprintf("%s:%d:%s", host, port, event),
		Host:      host,
		Port:      port,
		Event:     event,
		Severity:  sev,
		Timestamp: time.Now().UTC(),
		Labels:    make(map[string]string),
	}
}

// WithLabel attaches a key-value label to the envelope and returns the
// same pointer to allow chaining.
func (e *Envelope) WithLabel(key, value string) *Envelope {
	e.Labels[key] = value
	return e
}

// String returns a compact, human-readable representation of the envelope.
func (e *Envelope) String() string {
	return fmt.Sprintf("[%s] %s %s:%d (%s)",
		e.Severity, e.Event, e.Host, e.Port,
		e.Timestamp.Format(time.RFC3339))
}
