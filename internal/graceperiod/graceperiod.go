// Package graceperiod suppresses alerts for ports that appear and disappear
// within a short observation window, avoiding noise from transient connections.
package graceperiod

import (
	"errors"
	"sync"
	"time"
)

// Clock is a func that returns the current time; injectable for testing.
type Clock func() time.Time

// Tracker holds first-seen timestamps for ports and decides whether enough
// time has elapsed to consider a change "stable".
type Tracker struct {
	mu       sync.Mutex
	window   time.Duration
	clock    Clock
	firstSeen map[int]time.Time
}

// New returns a Tracker that requires a port to be observed for at least
// window before Allow returns true.
func New(window time.Duration, clock Clock) (*Tracker, error) {
	if window < 0 {
		return nil, errors.New("graceperiod: window must be non-negative")
	}
	if clock == nil {
		clock = time.Now
	}
	return &Tracker{
		window:    window,
		clock:     clock,
		firstSeen: make(map[int]time.Time),
	}, nil
}

// Observe records the first time a port is seen. Subsequent calls for the
// same port are no-ops until the port is forgotten via Forget.
func (t *Tracker) Observe(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.firstSeen[port]; !ok {
		t.firstSeen[port] = t.clock()
	}
}

// Allow returns true when the port has been observed for at least the
// configured window. A zero window always returns true.
func (t *Tracker) Allow(port int) bool {
	if t.window == 0 {
		return true
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	first, ok := t.firstSeen[port]
	if !ok {
		return false
	}
	return t.clock().Sub(first) >= t.window
}

// Forget removes the recorded first-seen time for a port so the grace period
// restarts if the port reappears.
func (t *Tracker) Forget(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.firstSeen, port)
}
