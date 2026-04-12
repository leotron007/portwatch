// Package baseline manages the set of ports that are explicitly trusted by the
// operator. Ports in the baseline are excluded from unexpected-change alerts,
// allowing portwatch to focus only on ports that deviate from the known-good
// configuration.
//
// A baseline is persisted as a JSON file so that it survives daemon restarts.
// Concurrent access is safe; all mutating operations hold an exclusive lock
// before modifying in-memory state and flushing to disk.
//
// Typical usage:
//
//	b, err := baseline.New("/var/lib/portwatch/baseline.json")
//	if err != nil { ... }
//
//	if !b.Contains(port) {
//		// alert: unexpected open port
//	}
package baseline
