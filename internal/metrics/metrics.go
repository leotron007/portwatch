// Package metrics tracks runtime counters for portwatch scans.
package metrics

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"text/tabwriter"
	"time"
)

// Counters holds atomic scan-cycle statistics.
type Counters struct {
	Scans    atomic.Int64
	Opened   atomic.Int64
	Closed   atomic.Int64
	Alerts   atomic.Int64
	Errors   atomic.Int64
	Started  time.Time
}

// New returns a Counters instance with the start time set to now.
func New() *Counters {
	return &Counters{Started: time.Now()}
}

// RecordScan increments the scan counter.
func (c *Counters) RecordScan() { c.Scans.Add(1) }

// RecordOpened increments the opened-ports counter by n.
func (c *Counters) RecordOpened(n int) { c.Opened.Add(int64(n)) }

// RecordClosed increments the closed-ports counter by n.
func (c *Counters) RecordClosed(n int) { c.Closed.Add(int64(n)) }

// RecordAlert increments the alert counter.
func (c *Counters) RecordAlert() { c.Alerts.Add(1) }

// RecordError increments the error counter.
func (c *Counters) RecordError() { c.Errors.Add(1) }

// Uptime returns the duration since the Counters were created.
func (c *Counters) Uptime() time.Duration { return time.Since(c.Started) }

// Write prints a human-readable summary to w (defaults to os.Stdout).
func (c *Counters) Write(w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "uptime\t%s\n", c.Uptime().Round(time.Second))
	fmt.Fprintf(tw, "scans\t%d\n", c.Scans.Load())
	fmt.Fprintf(tw, "opened\t%d\n", c.Opened.Load())
	fmt.Fprintf(tw, "closed\t%d\n", c.Closed.Load())
	fmt.Fprintf(tw, "alerts\t%d\n", c.Alerts.Load())
	fmt.Fprintf(tw, "errors\t%d\n", c.Errors.Load())
	return tw.Flush()
}
