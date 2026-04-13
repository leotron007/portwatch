// Package escalation provides a mechanism to escalate alerts when a
// condition persists beyond a configurable threshold of occurrences.
package escalation

import (
	"fmt"
	"sync"
	"time"
)

// Level represents the severity of an escalated alert.
type Level int

const (
	LevelNormal  Level = iota // below threshold
	LevelWarning              // at or above warning threshold
	LevelCritical             // at or above critical threshold
)

// String returns a human-readable label for the level.
func (l Level) String() string {
	switch l {
	case LevelWarning:
		return "warning"
	case LevelCritical:
		return "critical"
	default:
		return "normal"
	}
}

// Config holds thresholds for escalation.
type Config struct {
	WarningAfter  int           // number of occurrences before warning
	CriticalAfter int           // number of occurrences before critical
	DecayWindow   time.Duration // reset counter if key is idle this long
}

// entry tracks hit count and last seen time for a key.
type entry struct {
	count    int
	lastSeen time.Time
}

// Tracker counts repeated events per key and returns the current escalation level.
type Tracker struct {
	cfg   Config
	mu    sync.Mutex
	state map[string]*entry
	now   func() time.Time
}

// New creates a Tracker with the given Config.
// Returns an error if thresholds are non-positive or inconsistent.
func New(cfg Config) (*Tracker, error) {
	if cfg.WarningAfter <= 0 {
		return nil, fmt.Errorf("escalation: WarningAfter must be > 0, got %d", cfg.WarningAfter)
	}
	if cfg.CriticalAfter <= 0 {
		return nil, fmt.Errorf("escalation: CriticalAfter must be > 0, got %d", cfg.CriticalAfter)
	}
	if cfg.CriticalAfter <= cfg.WarningAfter {
		return nil, fmt.Errorf("escalation: CriticalAfter (%d) must be > WarningAfter (%d)", cfg.CriticalAfter, cfg.WarningAfter)
	}
	return &Tracker{
		cfg:   cfg,
		state: make(map[string]*entry),
		now:   time.Now,
	}, nil
}

// Record increments the hit counter for key and returns the resulting Level.
// If DecayWindow > 0 and the key has been idle longer than that window, the
// counter is reset before incrementing.
func (t *Tracker) Record(key string) Level {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	e, ok := t.state[key]
	if !ok {
		e = &entry{}
		t.state[key] = e
	}

	if t.cfg.DecayWindow > 0 && ok && now.Sub(e.lastSeen) > t.cfg.DecayWindow {
		e.count = 0
	}

	e.count++
	e.lastSeen = now

	return t.level(e.count)
}

// Reset clears the counter for the given key.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.state, key)
}

func (t *Tracker) level(count int) Level {
	switch {
	case count >= t.cfg.CriticalAfter:
		return LevelCritical
	case count >= t.cfg.WarningAfter:
		return LevelWarning
	default:
		return LevelNormal
	}
}
