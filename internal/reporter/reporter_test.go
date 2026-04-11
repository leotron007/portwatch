package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/state"
)

func TestReport_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Report(state.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no port changes") {
		t.Errorf("expected 'no port changes' in output, got: %q", buf.String())
	}
}

func TestReport_TextFormat_ShowsOpenedAndClosed(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	diff := state.Diff{
		Opened: []int{8080, 9090},
		Closed: []int{3000},
	}
	if err := r.Report(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"OPENED", "8080", "9090", "CLOSED", "3000"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestReport_JSONFormat_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	diff := state.Diff{
		Opened: []int{443},
		Closed: []int{80},
	}
	if err := r.Report(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"timestamp", "opened", "closed", "443", "80"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output, got: %q", want, out)
		}
	}
}

func TestNew_NilWriterDefaultsToStdout(t *testing.T) {
	// Should not panic when out is nil.
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
