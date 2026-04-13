package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Reason    string    `json:"reason,omitempty"`
}

// Logger writes structured audit entries to an io.Writer.
type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

// New returns a new Logger writing to w.
// If w is nil, os.Stderr is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{out: w}
}

// Log records an audit entry.
func (l *Logger) Log(event string, port int, proto, reason string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Port:      port,
		Proto:     proto,
		Reason:    reason,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}

	_, err = fmt.Fprintf(l.out, "%s\n", b)
	return err
}
