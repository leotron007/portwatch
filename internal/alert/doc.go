// Package alert implements the alerting subsystem for portwatch.
//
// When the scanner detects a port state change that matches a rule,
// an [Event] is constructed and dispatched to one or more [Notifier]
// implementations.
//
// # Notifiers
//
// The package ships with [LogNotifier], which writes human-readable
// alert lines to any [io.Writer] (defaults to os.Stderr).  Additional
// notifiers (e.g. webhook, syslog) can be added by implementing the
// [Notifier] interface.
//
// # Event Levels
//
// Three severity levels are defined:
//
//   - [LevelInfo]  — informational state change (e.g. expected port opened)
//   - [LevelWarn]  — advisory; rule matched but action is "warn"
//   - [LevelAlert] — critical; rule matched with action "alert" or "block"
package alert
