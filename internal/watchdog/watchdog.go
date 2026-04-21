// Package watchdog provides a self-monitoring component that detects
// when the scan loop stalls or stops making progress, emitting an alert
// if no heartbeat is received within the configured deadline.
package watchdog

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Watchdog monitors a heartbeat channel and writes an alert when the
// deadline is exceeded between consecutive beats.
type Watchdog struct {
	deadline time.Duration
	writer   io.Writer
	clock    func() time.Time
}

// Option is a functional option for Watchdog.
type Option func(*Watchdog)

// WithWriter sets the output writer used for stall alerts.
func WithWriter(w io.Writer) Option {
	return func(wd *Watchdog) { wd.writer = w }
}

// WithClock overrides the time source (useful in tests).
func WithClock(fn func() time.Time) Option {
	return func(wd *Watchdog) { wd.clock = fn }
}

// New creates a Watchdog that fires if no beat arrives within deadline.
// deadline must be positive.
func New(deadline time.Duration, opts ...Option) (*Watchdog, error) {
	if deadline <= 0 {
		return nil, fmt.Errorf("watchdog: deadline must be positive, got %s", deadline)
	}
	wd := &Watchdog{
		deadline: deadline,
		writer:   os.Stderr,
		clock:    time.Now,
	}
	for _, o := range opts {
		o(wd)
	}
	return wd, nil
}

// Run watches beats until ctx is cancelled. Each value received on beats
// resets the deadline timer. If the timer fires, a stall message is written
// to the configured writer.
func (wd *Watchdog) Run(ctx context.Context, beats <-chan struct{}) {
	var mu sync.Mutex
	last := wd.clock()

	ticker := time.NewTicker(wd.deadline / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-beats:
			if !ok {
				return
			}
			mu.Lock()
			last = wd.clock()
			mu.Unlock()
		case <-ticker.C:
			mu.Lock()
			since := wd.clock().Sub(last)
			mu.Unlock()
			if since >= wd.deadline {
				fmt.Fprintf(wd.writer, "[watchdog] stall detected: no heartbeat for %s\n", since.Round(time.Second))
			}
		}
	}
}
