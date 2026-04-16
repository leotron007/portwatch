package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNew_SetsStartTime(t *testing.T) {
	before := time.Now()
	c := metrics.New()
	after := time.Now()
	if c.Started.Before(before) || c.Started.After(after) {
		t.Fatalf("Started %v not in [%v, %v]", c.Started, before, after)
	}
}

func TestCounters_Increments(t *testing.T) {
	c := metrics.New()
	c.RecordScan()
	c.RecordScan()
	c.RecordOpened(3)
	c.RecordClosed(1)
	c.RecordAlert()
	c.RecordError()

	if got := c.Scans.Load(); got != 2 {
		t.Errorf("Scans = %d, want 2", got)
	}
	if got := c.Opened.Load(); got != 3 {
		t.Errorf("Opened = %d, want 3", got)
	}
	if got := c.Closed.Load(); got != 1 {
		t.Errorf("Closed = %d, want 1", got)
	}
	if got := c.Alerts.Load(); got != 1 {
		t.Errorf("Alerts = %d, want 1", got)
	}
	if got := c.Errors.Load(); got != 1 {
		t.Errorf("Errors = %d, want 1", got)
	}
}

func TestUptime_Positive(t *testing.T) {
	c := metrics.New()
	time.Sleep(2 * time.Millisecond)
	if c.Uptime() <= 0 {
		t.Error("expected positive uptime")
	}
}

func TestWrite_ContainsExpectedFields(t *testing.T) {
	c := metrics.New()
	c.RecordScan()
	c.RecordOpened(2)
	c.RecordClosed(1)
	c.RecordAlert()
	c.RecordError()

	var buf bytes.Buffer
	if err := c.Write(&buf); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"uptime", "scans", "opened", "closed", "alerts", "errors"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing field %q\n%s", want, out)
		}
	}
}

func TestWrite_NilWriterDefaultsToStdout(t *testing.T) {
	c := metrics.New()
	// Should not panic.
	if err := c.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
