// Package snapshot provides point-in-time captures of open port sets.
//
// A Snapshot pairs a sorted, deduplicated list of ports with a content-
// addressable digest and the wall-clock time of capture.  Two snapshots
// can be compared with Equal, or diffed with Added/Removed to surface
// ports that appeared or disappeared between consecutive scans.
//
// Typical usage:
//
//	old, _ := snapshot.New(previousPorts, time.Now())
//	// … time passes, re-scan …
//	current, _ := snapshot.New(latestPorts, time.Now())
//	if !old.Equal(current) {
//		opened := old.Added(current)
//		closed := old.Removed(current)
//	}
package snapshot
