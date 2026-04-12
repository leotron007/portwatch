// Package throttle provides a rate-limiter for alert notifications,
// suppressing repeated alerts for the same port within a configurable
// cooldown window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks the last alert time per port and suppresses duplicate
// notifications that occur within the cooldown duration.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSeen map[uint16]time.Time
	now      func() time.Time // injectable for testing
}

// New creates a Throttle with the given cooldown window.
// A cooldown of zero disables suppression (every event passes).
func New(cooldown time.Duration) *Throttle {
	return &Throttle{
		cooldown: cooldown,
		lastSeen: make(map[uint16]time.Time),
		now:      time.Now,
	}
}

// Allow reports whether an alert for the given port should be forwarded.
// It returns true the first time a port is seen and again only after the
// cooldown window has elapsed since the previous allowed alert.
func (t *Throttle) Allow(port uint16) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()

	if t.cooldown == 0 {
		return true
	}

	if last, ok := t.lastSeen[port]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}

	t.lastSeen[port] = now
	return true
}

// Reset clears the recorded state for a specific port, allowing the next
// alert for that port to pass immediately.
func (t *Throttle) Reset(port uint16) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSeen, port)
}

// Flush clears all recorded state.
func (t *Throttle) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastSeen = make(map[uint16]time.Time)
}
