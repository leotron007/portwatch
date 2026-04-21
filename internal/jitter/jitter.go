// Package jitter provides a configurable interval jitter utility that adds
// randomised offsets to a base duration, preventing thundering-herd effects
// when multiple goroutines wake on the same schedule.
package jitter

import (
	"fmt"
	"math/rand"
	"time"
)

// Source is a function that returns a pseudo-random float64 in [0.0, 1.0).
// It is exposed so tests can inject a deterministic source.
type Source func() float64

// Jitter adds a random offset up to MaxFraction of the base duration.
type Jitter struct {
	base        time.Duration
	maxFraction float64
	src         Source
}

// New returns a Jitter that offsets base by up to maxFraction (e.g. 0.25 for
// ±25 %). maxFraction must be in the range (0, 1].
func New(base time.Duration, maxFraction float64) (*Jitter, error) {
	if base <= 0 {
		return nil, fmt.Errorf("jitter: base duration must be positive, got %s", base)
	}
	if maxFraction <= 0 || maxFraction > 1 {
		return nil, fmt.Errorf("jitter: maxFraction must be in (0, 1], got %f", maxFraction)
	}
	return &Jitter{
		base:        base,
		maxFraction: maxFraction,
		src:         rand.Float64,
	}, nil
}

// withSource replaces the random source; used in tests.
func (j *Jitter) withSource(src Source) *Jitter {
	j.src = src
	return j
}

// Next returns the base a random offset in
// [0, base*maxFraction).
func (j *Jitter) Next() time.Duration {
	offset := time.Duration(float j.maxFraction * j.src())
	return j.base + offset
}

// Sleep blocks for Next() duration.
func (j *Jitter) Sleep() {
	time.Sleep(j.Next())
}
