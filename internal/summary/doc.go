// Package summary provides periodic scan summary reporting for the portwatch
// daemon.
//
// A Reporter reads entries from a history.History store and writes a
// human-readable summary line covering a configurable time window. It is
// intended to be invoked on a schedule (e.g. every hour or day) so operators
// can quickly see how many port-open and port-close events occurred without
// reading the full history log.
//
// Usage:
//
//	h, _ := history.New("/var/lib/portwatch/history.json")
//	r, _ := summary.New(h, os.Stdout, 24*time.Hour)
//	r.Write()
package summary
