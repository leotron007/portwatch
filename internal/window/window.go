// Package window provides a sliding-window counter used to track
// how many events have occurred within a rolling time period.
package window

import (
	"errors"
	"sync"
	"time"
)

// Clock abstracts time so tests can inject a fixed clock.
type Clock func() time.Time

// Window is a thread-safe sliding-window counter.
type Window struct {
	mu       sync.Mutex
	size     time.Duration
	buckets  int
	counts   []int
	times    []time.Time
	clock    Clock
}

// New creates a Window that divides [size] into [buckets] sub-intervals.
// size must be positive and buckets must be at least 1.
func New(size time.Duration, buckets int, clock Clock) (*Window, error) {
	if size <= 0 {
		return nil, errors.New("window: size must be positive")
	}
	if buckets < 1 {
		return nil, errors.New("window: buckets must be at least 1")
	}
	if clock == nil {
		clock = time.Now
	}
	return &Window{
		size:    size,
		buckets: buckets,
		counts:  make([]int, buckets),
		times:   make([]time.Time, buckets),
		clock:   clock,
	}, nil
}

// Add records n events at the current time.
func (w *Window) Add(n int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := w.clock()
	idx := w.bucketIndex(now)
	if w.times[idx].IsZero() || now.Sub(w.times[idx]) >= w.bucketSize() {
		w.counts[idx] = 0
		w.times[idx] = now
	}
	w.counts[idx] += n
}

// Count returns the total number of events recorded within the window.
func (w *Window) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := w.clock()
	total := 0
	for i := 0; i < w.buckets; i++ {
		if !w.times[i].IsZero() && now.Sub(w.times[i]) < w.size {
			total += w.counts[i]
		}
	}
	return total
}

// Reset clears all bucket data.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.counts = make([]int, w.buckets)
	w.times = make([]time.Time, w.buckets)
}

func (w *Window) bucketSize() time.Duration {
	return w.size / time.Duration(w.buckets)
}

func (w *Window) bucketIndex(t time.Time) int {
	bs := int64(w.bucketSize())
	if bs == 0 {
		return 0
	}
	return int((t.UnixNano()/bs)%int64(w.buckets)+int64(w.buckets)) % w.buckets
}
