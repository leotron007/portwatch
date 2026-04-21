// Package sampler provides adaptive scan interval adjustment based on
// observed port-change frequency. When changes are detected frequently the
// sampler shortens the interval; when the environment is stable it backs off
// toward a configurable maximum to reduce CPU and network load.
package sampler

import (
	"errors"
	"sync"
	"time"
)

// Sampler tracks recent change events and recommends a next scan interval.
type Sampler struct {
	mu      sync.Mutex
	min     time.Duration
	max     time.Duration
	current time.Duration
	// step is the factor applied on each stable cycle to grow the interval.
	step    float64
	events  int // change events observed in the current window
}

// New creates a Sampler with the given bounds and back-off step (e.g. 1.5).
// step must be > 1.0; min must be > 0 and <= max.
func New(min, max time.Duration, step float64) (*Sampler, error) {
	if min <= 0 {
		return nil, errors.New("sampler: min must be positive")
	}
	if max < min {
		return nil, errors.New("sampler: max must be >= min")
	}
	if step <= 1.0 {
		return nil, errors.New("sampler: step must be > 1.0")
	}
	return &Sampler{
		min:     min,
		max:     max,
		current: min,
		step:    step,
	}, nil
}

// Record signals that one or more port changes were observed in the last cycle.
// Calling Record resets the interval to the minimum.
func (s *Sampler) Record(changes int) {
	if changes <= 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events += changes
	s.current = s.min
}

// Advance should be called after each quiet (no-change) scan cycle. It grows
// the recommended interval by the configured step factor up to the maximum.
func (s *Sampler) Advance() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = 0
	next := time.Duration(float64(s.current) * s.step)
	if next > s.max {
		next = s.max
	}
	s.current = next
}

// Interval returns the currently recommended scan interval.
func (s *Sampler) Interval() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

// Reset restores the interval to the minimum.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = s.min
	s.events = 0
}
