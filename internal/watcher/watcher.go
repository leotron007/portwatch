// Package watcher ties together the scanner, filter, state, rules, alert,
// and reporter packages into a single scan-and-notify cycle.
package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Watcher performs periodic port scans and emits alerts on changes.
type Watcher struct {
	scanner  *scanner.Scanner
	filter   *filter.Filter
	state    *state.State
	rules    []rules.Rule
	notifier alert.Notifier
	reporter *reporter.Reporter
	interval time.Duration
}

// Config holds all dependencies needed to construct a Watcher.
type Config struct {
	Scanner  *scanner.Scanner
	Filter   *filter.Filter
	State    *state.State
	Rules    []rules.Rule
	Notifier alert.Notifier
	Reporter *reporter.Reporter
	Interval time.Duration
}

// New creates a Watcher from the provided Config.
// Returns an error if any required field is nil.
func New(cfg Config) (*Watcher, error) {
	if cfg.Scanner == nil {
		return nil, fmt.Errorf("watcher: scanner is required")
	}
	if cfg.Filter == nil {
		return nil, fmt.Errorf("watcher: filter is required")
	}
	if cfg.State == nil {
		return nil, fmt.Errorf("watcher: state is required")
	}
	if cfg.Notifier == nil {
		return nil, fmt.Errorf("watcher: notifier is required")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 30 * time.Second
	}
	return &Watcher{
		scanner:  cfg.Scanner,
		filter:   cfg.Filter,
		state:    cfg.State,
		rules:    cfg.Rules,
		notifier: cfg.Notifier,
		reporter: cfg.Reporter,
		interval: cfg.Interval,
	}, nil
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	if err := w.tick(ctx); err != nil {
		return err
	}
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(ctx); err != nil {
				return err
			}
		}
	}
}

// tick performs one scan cycle.
func (w *Watcher) tick(ctx context.Context) error {
	ports, err := w.scanner.Scan(ctx)
	if err != nil {
		return fmt.Errorf("watcher: scan failed: %w", err)
	}

	var allowed []int
	for _, p := range ports {
		if w.filter.Allow(p) {
			allowed = append(allowed, p)
		}
	}

	diff := w.state.Compare(allowed)

	for _, ev := range rules.Evaluate(w.rules, diff) {
		w.notifier.Notify(ev)
	}

	if w.reporter != nil {
		w.reporter.Report(diff)
	}

	return w.state.Save(allowed)
}
