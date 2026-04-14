// Package heartbeat provides a periodic tick mechanism that emits
// a signal at a configurable interval, allowing the daemon to confirm
// it is still running and trigger scheduled scan cycles.
package heartbeat

import (
	"context"
	"time"
)

// Beat is emitted on every heartbeat tick.
type Beat struct {
	// At is the time the beat was generated.
	At time.Time
	// Seq is the monotonically increasing beat counter (1-based).
	Seq uint64
}

// Ticker emits periodic Beat values on a channel.
type Ticker struct {
	interval time.Duration
	c        chan Beat
}

// New creates a new Ticker with the given interval.
// It returns an error if interval is less than or equal to zero.
func New(interval time.Duration) (*Ticker, error) {
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}
	return &Ticker{
		interval: interval,
		c:        make(chan Beat, 1),
	}, nil
}

// C returns the read-only channel on which beats are delivered.
func (t *Ticker) C() <-chan Beat {
	return t.c
}

// Run starts the ticker, emitting a Beat on every interval tick.
// It blocks until ctx is cancelled, then closes the channel.
func (t *Ticker) Run(ctx context.Context) {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()
	defer close(t.c)

	var seq uint64
	for {
		select {
		case <-ctx.Done():
			return
		case at := <-ticker.C:
			seq++
			select {
			case t.c <- Beat{At: at, Seq: seq}:
			default:
				// Drop the beat if the consumer is not keeping up.
			}
		}
	}
}

// ErrInvalidInterval is returned when a non-positive interval is provided.
var ErrInvalidInterval = errInvalidInterval("heartbeat: interval must be greater than zero")

type errInvalidInterval string

func (e errInvalidInterval) Error() string { return string(e) }
