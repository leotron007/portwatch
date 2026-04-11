// Package daemon wires together the scanner, state, rules, and alert
// subsystems and runs the port-watch polling loop.
package daemon

import (
	"context"
	"log"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/rules"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/state"
)

// Config holds runtime configuration for the daemon.
type Config struct {
	Host      string
	PortRange string
	Interval  time.Duration
	StatePath string
	Rules     []rules.Rule
	Notifier  alert.Notifier
}

// Daemon polls open ports on a fixed interval and fires alerts on changes.
type Daemon struct {
	cfg     Config
	scanner *scanner.Scanner
	state   *state.State
}

// New creates a Daemon from the supplied Config.
func New(cfg Config) (*Daemon, error) {
	s, err := scanner.New(cfg.Host, cfg.PortRange)
	if err != nil {
		return nil, err
	}
	st, err := state.New(cfg.StatePath)
	if err != nil {
		return nil, err
	}
	return &Daemon{cfg: cfg, scanner: s, state: st}, nil
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	log.Printf("portwatch daemon started (interval=%s)", d.cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch daemon stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := d.tick(); err != nil {
				log.Printf("tick error: %v", err)
			}
		}
	}
}

func (d *Daemon) tick() error {
	ports, err := d.scanner.Scan()
	if err != nil {
		return err
	}

	opened, closed := state.Compare(d.state.Ports(), ports)

	for _, p := range opened {
		ev := alert.Event{Port: p, Change: alert.Opened}
		if action, matched := rules.Evaluate(d.cfg.Rules, ev); matched {
			ev.Action = action
		}
		_ = d.cfg.Notifier.Notify(ev)
	}
	for _, p := range closed {
		ev := alert.Event{Port: p, Change: alert.Closed}
		if action, matched := rules.Evaluate(d.cfg.Rules, ev); matched {
			ev.Action = action
		}
		_ = d.cfg.Notifier.Notify(ev)
	}

	return d.state.Save(ports)
}
