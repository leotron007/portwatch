// Package metrics provides lightweight atomic counters for tracking
// portwatch daemon runtime statistics such as the number of scans
// performed, ports opened or closed, alerts fired, and errors
// encountered during a session.
//
// Usage:
//
//	c := metrics.New()
//	c.RecordScan()
//	c.RecordOpened(2)
//	c.Write(os.Stdout)
package metrics
