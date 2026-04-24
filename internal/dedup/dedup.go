// Package dedup provides event deduplication based on a content fingerprint
// and a configurable time window. Duplicate events observed within the window
// are suppressed so that downstream channels are not flooded with identical
// alerts.
package dedup

import (
	"errors"
	"sync"
	"time"
)

// Clock abstracts time so tests can inject a fixed instant.
type Clock func() time.Time

// Deduplicator tracks the last time each unique key was seen and suppresses
// repeated occurrences that fall within the configured window.
type Deduplicator struct {
	mu     sync.Mutex
	seen   map[string]time.Time
	window time.Duration
	clock  Clock
}

// New creates a Deduplicator with the given deduplication window.
// A zero window disables deduplication (every event is allowed).
// A negative window returns an error.
func New(window time.Duration, clock Clock) (*Deduplicator, error) {
	if window < 0 {
		return nil, errors.New("dedup: window must be non-negative")
	}
	if clock == nil {
		clock = time.Now
	}
	return &Deduplicator{
		seen:   make(map[string]time.Time),
		window: window,
		clock:  clock,
	}, nil
}

// Allow returns true when the event identified by key should be forwarded.
// It returns false when an identical key was already seen within the window.
func (d *Deduplicator) Allow(key string) bool {
	if d.window == 0 {
		return true
	}
	now := d.clock()
	d.mu.Lock()
	defer d.mu.Unlock()
	if last, ok := d.seen[key]; ok && now.Sub(last) < d.window {
		return false
	}
	d.seen[key] = now
	return true
}

// Flush removes all tracked keys, resetting the deduplication state.
func (d *Deduplicator) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}

// Len returns the number of keys currently tracked.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
