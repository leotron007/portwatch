// Package plateau detects when a metric has plateaued — remaining
// within a narrow band for a sustained number of consecutive observations.
package plateau

import (
	"errors"
	"sync"
)

// Detector tracks consecutive observations and reports when a value
// has remained within [baseline-tolerance, baseline+tolerance] for
// at least minRuns successive calls to Record.
type Detector struct {
	mu        sync.Mutex
	tolerance float64
	minRuns   int
	baseline  float64
	runs      int
	seeded    bool
}

// New creates a Detector. tolerance is the maximum absolute deviation
// from the first observed value that still counts as "flat".
// minRuns is the number of consecutive in-band observations required
// before Plateaued returns true.
func New(tolerance float64, minRuns int) (*Detector, error) {
	if tolerance < 0 {
		return nil, errors.New("plateau: tolerance must be non-negative")
	}
	if minRuns < 1 {
		return nil, errors.New("plateau: minRuns must be at least 1")
	}
	return &Detector{tolerance: tolerance, minRuns: minRuns}, nil
}

// Record submits the next observation. The first call seeds the baseline.
// Subsequent calls within tolerance increment the run counter; an
// out-of-band value resets both the baseline and the counter.
func (d *Detector) Record(v float64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.seeded {
		d.baseline = v
		d.runs = 1
		d.seeded = true
		return
	}

	diff := v - d.baseline
	if diff < 0 {
		diff = -diff
	}
	if diff <= d.tolerance {
		d.runs++
	} else {
		d.baseline = v
		d.runs = 1
	}
}

// Plateaued returns true when the required number of consecutive
// in-band observations has been reached.
func (d *Detector) Plateaued() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.seeded && d.runs >= d.minRuns
}

// Reset clears all state so the detector can be reused.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.runs = 0
	d.seeded = false
	d.baseline = 0
}
