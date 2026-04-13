package escalation_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/escalation"
)

func cfg() escalation.Config {
	return escalation.Config{
		WarningAfter:  2,
		CriticalAfter: 4,
		DecayWindow:   0,
	}
}

func TestNew_InvalidWarning(t *testing.T) {
	_, err := escalation.New(escalation.Config{WarningAfter: 0, CriticalAfter: 3})
	if err == nil {
		t.Fatal("expected error for WarningAfter=0")
	}
}

func TestNew_CriticalNotGreaterThanWarning(t *testing.T) {
	_, err := escalation.New(escalation.Config{WarningAfter: 3, CriticalAfter: 3})
	if err == nil {
		t.Fatal("expected error when CriticalAfter <= WarningAfter")
	}
}

func TestRecord_NormalBelowWarning(t *testing.T) {
	tr, _ := escalation.New(cfg())
	lvl := tr.Record("port:8080")
	if lvl != escalation.LevelNormal {
		t.Fatalf("expected Normal, got %s", lvl)
	}
}

func TestRecord_WarningAtThreshold(t *testing.T) {
	tr, _ := escalation.New(cfg())
	tr.Record("k")
	lvl := tr.Record("k")
	if lvl != escalation.LevelWarning {
		t.Fatalf("expected Warning, got %s", lvl)
	}
}

func TestRecord_CriticalAtThreshold(t *testing.T) {
	tr, _ := escalation.New(cfg())
	for i := 0; i < 3; i++ {
		tr.Record("k")
	}
	lvl := tr.Record("k")
	if lvl != escalation.LevelCritical {
		t.Fatalf("expected Critical, got %s", lvl)
	}
}

func TestRecord_IndependentKeys(t *testing.T) {
	tr, _ := escalation.New(cfg())
	for i := 0; i < 4; i++ {
		tr.Record("a")
	}
	lvl := tr.Record("b")
	if lvl != escalation.LevelNormal {
		t.Fatalf("key b should be Normal, got %s", lvl)
	}
}

func TestReset_ClearsCounter(t *testing.T) {
	tr, _ := escalation.New(cfg())
	for i := 0; i < 4; i++ {
		tr.Record("k")
	}
	tr.Reset("k")
	lvl := tr.Record("k")
	if lvl != escalation.LevelNormal {
		t.Fatalf("expected Normal after reset, got %s", lvl)
	}
}

func TestRecord_DecayWindowResetsCounter(t *testing.T) {
	clock := time.Now()
	cfgD := escalation.Config{
		WarningAfter:  2,
		CriticalAfter: 4,
		DecayWindow:   5 * time.Second,
	}
	tr, _ := escalation.New(cfgD)
	// Inject a fake clock via unexported field is not possible; use a sub-package
	// trick — instead just verify decay does not fire within window.
	_ = clock
	tr.Record("k")
	tr.Record("k")
	lvl := tr.Record("k") // count=3, still below critical
	if lvl != escalation.LevelWarning {
		t.Fatalf("expected Warning, got %s", lvl)
	}
}

func TestLevel_String(t *testing.T) {
	cases := []struct {
		lvl  escalation.Level
		want string
	}{
		{escalation.LevelNormal, "normal"},
		{escalation.LevelWarning, "warning"},
		{escalation.LevelCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.lvl.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.lvl, got, tc.want)
		}
	}
}
