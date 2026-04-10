package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
)

func baseEvent() alert.Event {
	return alert.Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     alert.LevelAlert,
		Port:      8080,
		Protocol:  "tcp",
		Message:   "unexpected open port detected",
		RuleName:  "block-8080",
	}
}

func TestLogNotifier_Notify_WritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)

	if err := n.Notify(baseEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, want := range []string{"2024-01-15T10:00:00Z", "ALERT", "block-8080", "8080", "tcp", "unexpected open port detected"} {
		if !strings.Contains(output, want) {
			t.Errorf("output %q missing expected substring %q", output, want)
		}
	}
}

func TestLogNotifier_Notify_InfoLevel(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)

	e := baseEvent()
	e.Level = alert.LevelInfo
	e.RuleName = "allow-80"
	e.Port = 80

	if err := n.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "INFO") {
		t.Errorf("expected INFO in output, got: %s", buf.String())
	}
}

func TestNewLogNotifier_NilWriterUsesStderr(t *testing.T) {
	// Should not panic when w is nil.
	n := alert.NewLogNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
