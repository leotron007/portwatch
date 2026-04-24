// Package stagger spreads a set of tasks over a time window to avoid
// thundering-herd effects when many goroutines wake simultaneously.
package stagger

import (
	"errors"
	"time"
)

// Clock abstracts time so tests can run without real sleeps.
type Clock func() time.Time

// Stagger holds configuration for distributing work across a window.
type Stagger struct {
	count  int
	window time.Duration
	clock  Clock
}

// New returns a Stagger that will spread count slots across window.
// count must be >= 1 and window must be positive.
func New(count int, window time.Duration, clock Clock) (*Stagger, error) {
	if count < 1 {
		return nil, errors.New("stagger: count must be >= 1")
	}
	if window <= 0 {
		return nil, errors.New("stagger: window must be positive")
	}
	if clock == nil {
		clock = time.Now
	}
	return &Stagger{count: count, window: window, clock: clock}, nil
}

// Delay returns the duration a given slot index should wait before
// executing. Slot indices are zero-based. An index outside [0, count)
// is clamped to the nearest boundary.
func (s *Stagger) Delay(slot int) time.Duration {
	if slot < 0 {
		slot = 0
	}
	if slot >= s.count {
		slot = s.count - 1
	}
	if s.count == 1 {
		return 0
	}
	step := s.window / time.Duration(s.count)
	return step * time.Duration(slot)
}

// Slots returns a slice of absolute wake-up times, one per slot,
// anchored to the current time returned by the configured clock.
func (s *Stagger) Slots() []time.Time {
	now := s.clock()
	slots := make([]time.Time, s.count)
	for i := 0; i < s.count; i++ {
		slots[i] = now.Add(s.Delay(i))
	}
	return slots
}

// Count returns the number of slots configured.
func (s *Stagger) Count() int { return s.count }

// Window returns the total spread duration.
func (s *Stagger) Window() time.Duration { return s.window }
