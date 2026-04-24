// Package flap detects ports that rapidly open and close (flapping).
// A port is considered flapping when it changes state more than a
// configured number of times within a sliding time window.
package flap

import (
	"errors"
	"sync"
	"time"
)

// Detector tracks state changes per port and reports flapping.
type Detector struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	clock     func() time.Time
	events    map[uint16][]time.Time
}

// New returns a Detector that considers a port flapping when it
// changes state at least threshold times within window.
func New(window time.Duration, threshold int, clock func() time.Time) (*Detector, error) {
	if window <= 0 {
		return nil, errors.New("flap: window must be positive")
	}
	if threshold < 2 {
		return nil, errors.New("flap: threshold must be at least 2")
	}
	if clock == nil {
		clock = time.Now
	}
	return &Detector{
		window:    window,
		threshold: threshold,
		clock:     clock,
		events:    make(map[uint16][]time.Time),
	}, nil
}

// Record registers a state-change event for port and returns true if
// the port is currently considered flapping.
func (d *Detector) Record(port uint16) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	cutoff := now.Add(-d.window)

	ts := d.events[port]
	// Prune events outside the window.
	valid := ts[:0]
	for _, t := range ts {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	valid = append(valid, now)
	d.events[port] = valid

	return len(valid) >= d.threshold
}

// Reset clears the event history for port.
func (d *Detector) Reset(port uint16) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.events, port)
}

// Flapping returns all ports currently considered flapping.
func (d *Detector) Flapping() []uint16 {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	cutoff := now.Add(-d.window)
	var out []uint16
	for port, ts := range d.events {
		count := 0
		for _, t := range ts {
			if t.After(cutoff) {
				count++
			}
		}
		if count >= d.threshold {
			out = append(out, port)
		}
	}
	return out
}
