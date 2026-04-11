// Package alert provides alerting mechanisms for portwatch,
// notifying users when port state changes violate configured rules.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port change that triggered an alert.
type Event struct {
	Timestamp time.Time
	Level     Level
	Port      int
	Protocol  string
	Message   string
	RuleName  string
}

// String returns a human-readable representation of the event.
func (e Event) String() string {
	return fmt.Sprintf(
		"%s [%s] rule=%q port=%d/%s — %s",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.RuleName,
		e.Port,
		e.Protocol,
		e.Message,
	)
}

// Notifier sends alert events to a destination.
type Notifier interface {
	Notify(e Event) error
}

// LogNotifier writes alert events as formatted lines to an io.Writer.
type LogNotifier struct {
	Writer io.Writer
}

// NewLogNotifier returns a LogNotifier that writes to w.
// If w is nil, os.Stderr is used.
func NewLogNotifier(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stderr
	}
	return &LogNotifier{Writer: w}
}

// Notify formats and writes the event to the underlying writer.
func (n *LogNotifier) Notify(e Event) error {
	_, err := fmt.Fprintln(n.Writer, e.String())
	return err
}
