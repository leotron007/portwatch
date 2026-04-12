// Package ratelimit provides a token-bucket style rate limiter that
// restricts how many alert events can be dispatched within a sliding
// time window.  It is safe for concurrent use.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Clock abstracts time so tests can inject a deterministic source.
type Clock func() time.Time

// Limiter tracks per-key event counts within a rolling window.
type Limiter struct {
	mu      sync.Mutex
	window  time.Duration
	max     int
	clock   Clock
	buckets map[string][]time.Time
}

// New creates a Limiter that allows at most max events per key inside
// the given window duration.  max must be >= 1 and window > 0.
func New(window time.Duration, max int, clock Clock) (*Limiter, error) {
	if window <= 0 {
		return nil, fmt.Errorf("ratelimit: window must be positive, got %s", window)
	}
	if max < 1 {
		return nil, fmt.Errorf("ratelimit: max must be >= 1, got %d", max)
	}
	if clock == nil {
		clock = time.Now
	}
	return &Limiter{
		window:  window,
		max:     max,
		clock:   clock,
		buckets: make(map[string][]time.Time),
	}, nil
}

// Allow reports whether a new event for key is permitted.  It prunes
// timestamps that have fallen outside the current window before
// deciding, so memory usage stays bounded.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	cutoff := now.Add(-l.window)

	times := l.buckets[key]
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.max {
		l.buckets[key] = valid
		return false
	}

	l.buckets[key] = append(valid, now)
	return true
}

// Reset clears all recorded events for key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}
