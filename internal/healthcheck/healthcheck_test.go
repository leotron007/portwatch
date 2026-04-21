package healthcheck_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/healthcheck"
)

func TestNew_NilWriterDefaultsToStdout(t *testing.T) {
	r := healthcheck.New(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestRun_AllHealthy(t *testing.T) {
	r := healthcheck.New(nil)
	r.Register(healthcheck.Check{Name: "scanner", Probe: func() error { return nil }})
	r.Register(healthcheck.Check{Name: "state", Probe: func() error { return nil }})

	results, status := r.Run()
	if status != healthcheck.StatusOK {
		t.Fatalf("expected ok, got %s", status)
	}
	for _, res := range results {
		if !res.Healthy {
			t.Errorf("expected %s to be healthy", res.Name)
		}
	}
}

func TestRun_NonCriticalFailureDegraded(t *testing.T) {
	r := healthcheck.New(nil)
	r.Register(healthcheck.Check{
		Name:     "optional",
		Critical: false,
		Probe:    func() error { return errors.New("not available") },
	})

	_, status := r.Run()
	if status != healthcheck.StatusDegraded {
		t.Fatalf("expected degraded, got %s", status)
	}
}

func TestRun_CriticalFailureUnhealthy(t *testing.T) {
	r := healthcheck.New(nil)
	r.Register(healthcheck.Check{
		Name:     "scanner",
		Critical: true,
		Probe:    func() error { return errors.New("scanner down") },
	})

	_, status := r.Run()
	if status != healthcheck.StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", status)
	}
}

func TestWrite_ContainsCheckName(t *testing.T) {
	var buf bytes.Buffer
	r := healthcheck.New(&buf)
	r.Register(healthcheck.Check{Name: "state-file", Probe: func() error { return nil }})

	r.Write()

	if !strings.Contains(buf.String(), "state-file") {
		t.Errorf("expected output to contain check name, got: %s", buf.String())
	}
}

func TestWrite_FailedCheckContainsMessage(t *testing.T) {
	var buf bytes.Buffer
	r := healthcheck.New(&buf)
	r.Register(healthcheck.Check{
		Name:  "history",
		Probe: func() error { return errors.New("file locked") },
	})

	r.Write()

	if !strings.Contains(buf.String(), "file locked") {
		t.Errorf("expected error message in output, got: %s", buf.String())
	}
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		s    healthcheck.Status
		want string
	}{
		{healthcheck.StatusOK, "ok"},
		{healthcheck.StatusDegraded, "degraded"},
		{healthcheck.StatusUnhealthy, "unhealthy"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}
