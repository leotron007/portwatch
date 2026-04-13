// Package debounce provides a port-event debouncer that suppresses
// transient open/close flaps by requiring a port state to remain
// stable for a configurable duration before forwarding the event.
package debounce

import (
	"sync"
	"time"
)

// Clock is a minimal time abstraction to allow deterministic testing.
type Clock func() time.Time

// entry tracks the first time a (port, state) pair was observed.
type entry struct {
	since time.Time
	state bool // true = open, false = closed
}

// Debouncer holds per-port observations and emits an event only after
// the port has remained in the same state for at least Window.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	clock   Clock
	pending map[int]entry
}

// New creates a Debouncer with the given stability window.
// A zero or negative window means every event passes immediately.
func New(window time.Duration, clock Clock) *Debouncer {
	if clock == nil {
		clock = time.Now
	}
	return &Debouncer{
		window:  window,
		clock:   clock,
		pending: make(map[int]entry),
	}
}

// Allow returns true when the port has been continuously observed in
// state open for at least the configured window, or immediately if the
// window is zero. Calling Allow with a different state for the same
// port resets the observation timer.
func (d *Debouncer) Allow(port int, open bool) bool {
	if d.window <= 0 {
		return true
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	e, exists := d.pending[port]

	if !exists || e.state != open {
		// First observation or state changed — reset the timer.
		d.pending[port] = entry{since: now, state: open}
		return false
	}

	if now.Sub(e.since) >= d.window {
		delete(d.pending, port)
		return true
	}

	return false
}

// Reset clears any pending observation for the given port.
func (d *Debouncer) Reset(port int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.pending, port)
}

// Len returns the number of ports currently in the pending state.
func (d *Debouncer) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.pending)
}
