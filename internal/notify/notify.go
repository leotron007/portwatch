package notify

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Channel represents a notification delivery mechanism.
type Channel interface {
	Send(event alert.Event) error
}

// ExecChannel runs an external command, passing event details via environment variables.
type ExecChannel struct {
	Command string
	Args    []string
}

// Send executes the configured command with event metadata in the environment.
func (e *ExecChannel) Send(event alert.Event) error {
	if e.Command == "" {
		return fmt.Errorf("notify: exec command must not be empty")
	}
	cmd := exec.Command(e.Command, e.Args...)
	cmd.Env = append(os.Environ(),
		"PORTWATCH_PORT="+fmt.Sprintf("%d", event.Port),
		"PORTWATCH_ACTION="+string(event.Action),
		"PORTWATCH_RULE="+event.RuleName,
		"PORTWATCH_HOST="+event.Host,
		"PORTWATCH_TIME="+event.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("notify: exec %q failed: %w (output: %s)", e.Command, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// WebhookChannel sends a JSON POST to a URL using the provided writer for logging.
type WriterChannel struct {
	Writer io.Writer
	Prefix string
}

// NewWriterChannel returns a WriterChannel that writes to w.
func NewWriterChannel(w io.Writer, prefix string) *WriterChannel {
	if w == nil {
		w = os.Stderr
	}
	return &WriterChannel{Writer: w, Prefix: prefix}
}

// Send formats the event and writes it to the underlying writer.
func (w *WriterChannel) Send(event alert.Event) error {
	_, err := fmt.Fprintf(w.Writer, "%s port=%d action=%s rule=%q host=%s time=%s\n",
		w.Prefix,
		event.Port,
		string(event.Action),
		event.RuleName,
		event.Host,
		event.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	)
	return err
}
