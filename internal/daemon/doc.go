// Package daemon implements the portwatch polling loop.
//
// The Daemon ties together the scanner, persistent state, rule engine, and
// alert notifier into a single long-running process:
//
//  1. On every tick the scanner probes the configured host:port-range.
//  2. The result is diffed against the previous snapshot held in State.
//  3. Each changed port is evaluated against the loaded Rules.
//  4. A Notifier fires for every change (with the matched action attached
//     when a rule applies).
//  5. The new snapshot is persisted so the next tick has a fresh baseline.
//
// Typical usage:
//
//	d, err := daemon.New(cfg)
//	if err != nil { log.Fatal(err) }
//	if err := d.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//	    log.Fatal(err)
//	}
package daemon
