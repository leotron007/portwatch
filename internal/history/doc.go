// Package history records port-change events over time and persists them to a
// newline-delimited JSON file on disk.
//
// Usage:
//
//	h, err := history.New("/var/lib/portwatch/history.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Record an event whenever the watcher detects a change.
//	_ = h.Record(8080, "opened", "tcp")
//
//	// Retrieve all past events for reporting.
//	entries := h.Entries()
//
The file is created if it does not exist. Entries are appended in memory and
the full slice is re-serialised to disk on every call to Record, keeping the
implementation simple at the cost of write amplification for very large
histories.
package history
