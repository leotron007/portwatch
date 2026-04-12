// Package suppress provides a mechanism to suppress duplicate alerts
// for a given port+action combination within a configurable time window.
package suppress

import (
	"sync"
	"time"
)

// key uniquely identifies a suppressible event.
type key struct {
	port   int
	action string
}

// Suppressor tracks recently seen events and suppresses duplicates
// that occur within the configured window.
type Suppressor struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[key]time.Time
	now    func() time.Time
}

// New creates a Suppressor with the given deduplication window.
// A zero or negative window disables suppression (all events pass).
func New(window time.Duration) *Suppressor {
	return &Suppressor{
		window: window,
		seen:   make(map[key]time.Time),
		now:    time.Now,
	}
}

// Allow returns true if the event for the given port and action should
// be forwarded, or false if it is a duplicate within the active window.
func (s *Suppressor) Allow(port int, action string) bool {
	if s.window <= 0 {
		return true
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	k := key{port: port, action: action}

	if last, ok := s.seen[k]; ok && now.Sub(last) < s.window {
		return false
	}

	s.seen[k] = now
	return true
}

// Flush removes all entries whose window has expired, freeing memory
// for long-running daemon processes.
func (s *Suppressor) Flush() {
	if s.window <= 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	for k, last := range s.seen {
		if now.Sub(last) >= s.window {
			delete(s.seen, k)
		}
	}
}
