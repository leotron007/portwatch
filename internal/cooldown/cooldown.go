// Package cooldown provides per-port cooldown tracking to prevent
// alert fatigue by suppressing repeated notifications for the same port
// within a configurable duration.
package cooldown

import (
	"sync"
	"time"
)

// Clock allows time to be injected for testing.
type Clock func() time.Time

// Tracker records the last alert time for each port and decides whether
// a new alert should be allowed based on a per-port cooldown window.
type Tracker struct {
	mu       sync.Mutex
	last     map[int]time.Time
	cooldown time.Duration
	now      Clock
}

// New creates a Tracker with the given cooldown duration.
// A zero or negative cooldown means every event is allowed through.
func New(cooldown time.Duration) *Tracker {
	return &Tracker{
		last:     make(map[int]time.Time),
		cooldown: cooldown,
		now:      time.Now,
	}
}

// Allow returns true if the port has not been alerted within the cooldown
// window. If allowed, the last-seen time for the port is updated.
func (t *Tracker) Allow(port int) bool {
	if t.cooldown <= 0 {
		return true
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.last[port]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.last[port] = now
	return true
}

// Reset clears the cooldown record for a specific port.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, port)
}

// Len returns the number of ports currently tracked.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.last)
}
