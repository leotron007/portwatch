package backoff

import (
	"errors"
	"math"
	"sync"
	"time"
)

// Clock allows injecting a custom time source for testing.
type Clock func() time.Time

// Backoff implements exponential back-off with an optional jitter fraction
// and a configurable ceiling. It is safe for concurrent use.
type Backoff struct {
	mu       sync.Mutex
	attempts int
	base     time.Duration
	max      time.Duration
	factor   float64
	jitter   float64 // fraction of computed delay to randomise [0,1)
	clock    Clock
}

// New creates a Backoff. base is the initial delay, max is the ceiling,
// factor is the multiplier per attempt (≥1.0), and jitter is a fraction
// of the computed delay added as random noise (0 disables jitter).
func New(base, max time.Duration, factor, jitter float64) (*Backoff, error) {
	if base <= 0 {
		return nil, errors.New("backoff: base must be positive")
	}
	if max < base {
		return nil, errors.New("backoff: max must be >= base")
	}
	if factor < 1.0 {
		return nil, errors.New("backoff: factor must be >= 1.0")
	}
	if jitter < 0 || jitter >= 1 {
		return nil, errors.New("backoff: jitter must be in [0, 1)")
	}
	return &Backoff{
		base:   base,
		max:    max,
		factor: factor,
		jitter: jitter,
		clock:  time.Now,
	}, nil
}

// Next returns the delay for the current attempt and increments the counter.
func (b *Backoff) Next() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	delay := float64(b.base) * math.Pow(b.factor, float64(b.attempts))
	if delay > float64(b.max) {
		delay = float64(b.max)
	}
	if b.jitter > 0 {
		// deterministic-ish jitter seeded from wall time nanoseconds
		nano := float64(b.clock().UnixNano() % 1_000_000)
		frac := (nano / 1_000_000) * b.jitter
		delay += delay * frac
		if delay > float64(b.max) {
			delay = float64(b.max)
		}
	}
	b.attempts++
	return time.Duration(delay)
}

// Reset sets the attempt counter back to zero.
func (b *Backoff) Reset() {
	b.mu.Lock()
	b.attempts = 0
	b.mu.Unlock()
}

// Attempts returns the number of Next calls since the last Reset.
func (b *Backoff) Attempts() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.attempts
}
