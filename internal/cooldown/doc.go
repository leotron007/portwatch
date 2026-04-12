// Package cooldown implements per-port cooldown tracking for portwatch.
//
// A Tracker records the last alert timestamp for each port number and
// exposes an Allow method that returns false when the port has been seen
// within the configured cooldown window. This prevents repeated alerts
// for ports that oscillate rapidly between open and closed states.
//
// Basic usage:
//
//	tr := cooldown.New(30 * time.Second)
//	if tr.Allow(port) {
//		// send alert
//	}
//
// The cooldown window is configured per-daemon via the YAML key
// "cooldown_seconds" in the main configuration file.
package cooldown
