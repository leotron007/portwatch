// Package watcher orchestrates a complete port-monitoring cycle.
//
// A Watcher combines a [scanner.Scanner], [filter.Filter], [state.State],
// [rules.Rule] set, [alert.Notifier], and optional [reporter.Reporter] into a
// single, self-contained loop that:
//
//  1. Scans the configured host for open TCP ports.
//  2. Filters results according to the active [filter.Filter].
//  3. Compares the current snapshot against the persisted [state.State].
//  4. Evaluates [rules.Rule]s against the diff and forwards matching events to
//     the [alert.Notifier].
//  5. Optionally writes a human- or machine-readable summary via the
//     [reporter.Reporter].
//  6. Persists the new snapshot so the next cycle has an accurate baseline.
//
// Typical usage:
//
//	w, err := watcher.New(watcher.Config{...})
//	if err != nil { log.Fatal(err) }
//	if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package watcher
