package notify_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

var baseEvent = alert.Event{
	Port:      8080,
	Action:    alert.ActionOpened,
	RuleName:  "web-ports",
	Host:      "127.0.0.1",
	Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
}

func TestWriterChannel_Send_FormatsLine(t *testing.T) {
	var buf bytes.Buffer
	ch := notify.NewWriterChannel(&buf, "[ALERT]")
	if err := ch.Send(baseEvent); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"[ALERT]", "port=8080", "action=opened", "web-ports", "127.0.0.1"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in output %q", want, got)
		}
	}
}

func TestWriterChannel_NilWriterUsesStderr(t *testing.T) {
	ch := notify.NewWriterChannel(nil, "")
	if ch.Writer == nil {
		t.Fatal("expected non-nil writer when nil passed")
	}
}

func TestExecChannel_EmptyCommandReturnsError(t *testing.T) {
	ch := &notify.ExecChannel{Command: ""}
	err := ch.Send(baseEvent)
	if err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestExecChannel_InvalidCommandReturnsError(t *testing.T) {
	ch := &notify.ExecChannel{Command: "/nonexistent/binary"}
	err := ch.Send(baseEvent)
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
}

func TestExecChannel_ValidCommand(t *testing.T) {
	ch := &notify.ExecChannel{Command: "true"}
	if err := ch.Send(baseEvent); err != nil {
		t.Fatalf("unexpected error running 'true': %v", err)
	}
}
