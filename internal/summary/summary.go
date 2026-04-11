// Package summary provides periodic scan summary reporting for portwatch.
package summary

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/history"
)

// Reporter writes periodic summaries of port scan history to a writer.
type Reporter struct {
	history *history.History
	writer  io.Writer
	window  time.Duration
}

// New creates a new summary Reporter. If w is nil, os.Stdout is used.
func New(h *history.History, w io.Writer, window time.Duration) (*Reporter, error) {
	if h == nil {
		return nil, fmt.Errorf("summary: history must not be nil")
	}
	if w == nil {
		w = os.Stdout
	}
	if window <= 0 {
		window = 24 * time.Hour
	}
	return &Reporter{history: h, writer: w, window: window}, nil
}

// Write outputs a summary of events within the configured time window.
func (r *Reporter) Write() error {
	since := time.Now().Add(-r.window)
	entries := r.history.Since(since)

	opened := 0
	closed := 0
	for _, e := range entries {
		switch e.Kind {
		case "opened":
			opened++
		case "closed":
			closed++
		}
	}

	_, err := fmt.Fprintf(
		r.writer,
		"[summary] window=%s events=%d opened=%d closed=%d\n",
		r.window.String(),
		len(entries),
		opened,
		closed,
	)
	return err
}
