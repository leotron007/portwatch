package suppress

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	s := New(5 * time.Second)
	if !s.Allow(8080, "opened") {
		t.Fatal("expected first call to pass")
	}
}

func TestAllow_DuplicateWithinWindowBlocked(t *testing.T) {
	now := time.Now()
	s := New(10 * time.Second)
	s.now = fixedClock(now)

	s.Allow(443, "opened")

	if s.Allow(443, "opened") {
		t.Fatal("expected duplicate within window to be blocked")
	}
}

func TestAllow_PassesAfterWindowExpires(t *testing.T) {
	now := time.Now()
	s := New(5 * time.Second)
	s.now = fixedClock(now)

	s.Allow(80, "closed")

	s.now = fixedClock(now.Add(6 * time.Second))
	if !s.Allow(80, "closed") {
		t.Fatal("expected call after window expiry to pass")
	}
}

func TestAllow_ZeroWindowNeverBlocks(t *testing.T) {
	s := New(0)
	s.Allow(22, "opened")
	if !s.Allow(22, "opened") {
		t.Fatal("expected zero window to never suppress")
	}
}

func TestAllow_DifferentActionsAreIndependent(t *testing.T) {
	now := time.Now()
	s := New(10 * time.Second)
	s.now = fixedClock(now)

	s.Allow(8080, "opened")

	if !s.Allow(8080, "closed") {
		t.Fatal("expected different action on same port to pass")
	}
}

func TestFlush_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	s := New(5 * time.Second)
	s.now = fixedClock(now)

	s.Allow(9090, "opened")
	if len(s.seen) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(s.seen))
	}

	s.now = fixedClock(now.Add(10 * time.Second))
	s.Flush()

	if len(s.seen) != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", len(s.seen))
	}
}

func TestFlush_KeepsActiveEntries(t *testing.T) {
	now := time.Now()
	s := New(30 * time.Second)
	s.now = fixedClock(now)

	s.Allow(3000, "opened")

	s.now = fixedClock(now.Add(5 * time.Second))
	s.Flush()

	if len(s.seen) != 1 {
		t.Fatalf("expected entry to survive flush, got %d entries", len(s.seen))
	}
}
