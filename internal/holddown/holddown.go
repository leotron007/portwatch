// Package holddown implements a hold-down timer that requires a condition to
// persist for a minimum duration before it is considered stable and allowed
// through. This prevents transient port flaps from generating spurious alerts.
package holddown

import (
	"errors"
	"sync"
	"time"
)

// Clock is a functional dependency for time, allowing deterministic tests.
type Clock func() time.Time

// Tracker holds per-key hold-down state.
type Tracker struct {
	mu       sync.Mutex
	window   time.Duration
	clock    Clock
	first    map[string]time.Time
	lastSeen map[string]bool
}

// New creates a Tracker with the given hold-down window.
// window must be non-negative; a zero window causes every call to pass.
func New(window time.Duration, clock Clock) (*Tracker, error) {
	if window < 0 {
		return nil, errors.New("holddown: window must be non-negative")
	}
	if clock == nil {
		clock = time.Now
	}
	return &Tracker{
		window:   window,
		clock:    clock,
		first:    make(map[string]time.Time),
		lastSeen: make(map[string]bool),
	}, nil
}

// Allow reports whether key has been continuously observed for at least the
// configured hold-down window. The first observation starts the timer; if the
// key is absent on a subsequent call the timer resets.
func (t *Tracker) Allow(key string, present bool) bool {
	if t.window == 0 {
		return present
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()

	if !present {
		delete(t.first, key)
		delete(t.lastSeen, key)
		return false
	}

	if _, ok := t.first[key]; !ok {
		t.first[key] = now
		t.lastSeen[key] = true
		return false
	}

	t.lastSeen[key] = true
	return now.Sub(t.first[key]) >= t.window
}

// Reset clears the hold-down state for key.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.first, key)
	delete(t.lastSeen, key)
}
