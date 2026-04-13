// Package escalation tracks how many times a named event has been observed
// and returns a severity Level (Normal, Warning, Critical) based on
// configurable occurrence thresholds.
//
// A Tracker is safe for concurrent use. Each unique key maintains its own
// independent counter. An optional DecayWindow causes a key's counter to
// reset automatically when the key has not been seen for longer than the
// configured duration, preventing stale counts from triggering false
// escalations after a quiet period.
//
// Typical usage:
//
//	tr, err := escalation.New(escalation.Config{
//		WarningAfter:  3,
//		CriticalAfter: 6,
//		DecayWindow:   10 * time.Minute,
//	})
//	lvl := tr.Record(fmt.Sprintf("port:%d", port))
//	if lvl >= escalation.LevelWarning {
//		// send elevated alert
//	}
package escalation
