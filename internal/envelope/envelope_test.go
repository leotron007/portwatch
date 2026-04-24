package envelope_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/envelope"
)

func TestNew_FieldsPopulated(t *testing.T) {
	before := time.Now().UTC()
	e := envelope.New("127.0.0.1", 8080, envelope.EventOpened, envelope.SeverityInfo)
	after := time.Now().UTC()

	if e.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", e.Host)
	}
	if e.Port != 8080 {
		t.Errorf("expected port 8080, got %d", e.Port)
	}
	if e.Event != envelope.EventOpened {
		t.Errorf("expected event opened, got %s", e.Event)
	}
	if e.Severity != envelope.SeverityInfo {
		t.Errorf("expected severity info, got %v", e.Severity)
	}
	if e.Timestamp.Before(before) || e.Timestamp.After(after) {
		t.Errorf("timestamp %v outside expected range", e.Timestamp)
	}
}

func TestNew_IDIsNonEmpty(t *testing.T) {
	e := envelope.New("localhost", 443, envelope.EventClosed, envelope.SeverityWarning)
	if e.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestWithLabel_ChainedAndStored(t *testing.T) {
	e := envelope.New("localhost", 22, envelope.EventOpened, envelope.SeverityInfo).
		WithLabel("env", "prod").
		WithLabel("team", "ops")

	if e.Labels["env"] != "prod" {
		t.Errorf("expected label env=prod, got %s", e.Labels["env"])
	}
	if e.Labels["team"] != "ops" {
		t.Errorf("expected label team=ops, got %s", e.Labels["team"])
	}
}

func TestSeverity_String(t *testing.T) {
	cases := []struct {
		sev  envelope.Severity
		want string
	}{
		{envelope.SeverityInfo, "info"},
		{envelope.SeverityWarning, "warning"},
		{envelope.SeverityCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.sev.String(); got != tc.want {
			t.Errorf("Severity.String() = %q, want %q", got, tc.want)
		}
	}
}

func TestString_ContainsKeyParts(t *testing.T) {
	e := envelope.New("10.0.0.1", 3306, envelope.EventClosed, envelope.SeverityCritical)
	s := e.String()
	for _, part := range []string{"critical", "closed", "10.0.0.1", "3306"} {
		if !strings.Contains(s, part) {
			t.Errorf("String() missing %q: %s", part, s)
		}
	}
}
