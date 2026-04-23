// Package budget implements an error-budget tracker that accumulates
// failure events within a rolling time window and reports when the
// remaining budget drops below a configured threshold.
package budget

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Budget tracks error events against a fixed capacity and exposes the
// remaining budget as a fraction in [0, 1].
type Budget struct {
	mu       sync.Mutex
	capacity int
	window   time.Duration
	events   []time.Time
	clock    func() time.Time
}

// New creates a Budget with the given capacity (maximum allowed failures)
// and rolling window duration. capacity must be ≥ 1 and window must be > 0.
func New(capacity int, window time.Duration) (*Budget, error) {
	if capacity < 1 {
		return nil, errors.New("budget: capacity must be at least 1")
	}
	if window <= 0 {
		return nil, errors.New("budget: window must be positive")
	}
	return &Budget{
		capacity: capacity,
		window:   window,
		events:   make([]time.Time, 0, capacity),
		clock:    time.Now,
	}, nil
}

// Record registers one failure event at the current time.
func (b *Budget) Record() {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.clock()
	b.prune(now)
	b.events = append(b.events, now)
}

// Remaining returns the fraction of budget left in [0, 1].
// A value of 0 means the budget is exhausted.
func (b *Budget) Remaining() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.prune(b.clock())
	used := len(b.events)
	if used >= b.capacity {
		return 0
	}
	return float64(b.capacity-used) / float64(b.capacity)
}

// Exhausted reports whether the error budget has been fully consumed.
func (b *Budget) Exhausted() bool {
	return b.Remaining() == 0
}

// String returns a human-readable summary of the current budget state.
func (b *Budget) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.prune(b.clock())
	return fmt.Sprintf("budget: %d/%d used (%.0f%% remaining)",
		len(b.events), b.capacity,
		b.Remaining()*100)
}

// prune removes events that have fallen outside the rolling window.
// Must be called with b.mu held.
func (b *Budget) prune(now time.Time) {
	cutoff := now.Add(-b.window)
	i := 0
	for i < len(b.events) && b.events[i].Before(cutoff) {
		i++
	}
	b.events = b.events[i:]
}
