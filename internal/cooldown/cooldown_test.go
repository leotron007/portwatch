package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
)

type fixedClock struct{ t time.Time }

func (f *fixedClock) Now() time.Time { return f.t }

func newTracker(d time.Duration, clk *fixedClock) *cooldown.Tracker {
	tr := cooldown.New(d)
	tr.(*interface{ SetClock(cooldown.Clock) }) // compile-time unused; patched below
	_ = clk
	return tr
}

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	tr := cooldown.New(5 * time.Second)
	if !tr.Allow(8080) {
		t.Fatal("expected first call to pass")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	tr := cooldown.NewWithClock(5*time.Second, clk.Now)
	tr.Allow(8080)
	if tr.Allow(8080) {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_PassesAfterCooldownExpires(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	tr := cooldown.NewWithClock(5*time.Second, clk.Now)
	tr.Allow(8080)
	clk.t = clk.t.Add(6 * time.Second)
	if !tr.Allow(8080) {
		t.Fatal("expected call after cooldown to pass")
	}
}

func TestAllow_ZeroCooldownNeverBlocks(t *testing.T) {
	tr := cooldown.New(0)
	for i := 0; i < 5; i++ {
		if !tr.Allow(9090) {
			t.Fatalf("expected zero cooldown to always pass, failed on iteration %d", i)
		}
	}
}

func TestReset_ClearsPort(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	tr := cooldown.NewWithClock(10*time.Second, clk.Now)
	tr.Allow(443)
	tr.Reset(443)
	if !tr.Allow(443) {
		t.Fatal("expected allow after reset")
	}
}

func TestLen_TracksActiveEntries(t *testing.T) {
	tr := cooldown.New(10 * time.Second)
	tr.Allow(80)
	tr.Allow(443)
	tr.Allow(8080)
	if tr.Len() != 3 {
		t.Fatalf("expected 3 tracked ports, got %d", tr.Len())
	}
	tr.Reset(443)
	if tr.Len() != 2 {
		t.Fatalf("expected 2 tracked ports after reset, got %d", tr.Len())
	}
}
