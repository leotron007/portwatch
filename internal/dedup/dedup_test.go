package dedup

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) Clock {
	return func() time.Time { return t }
}

func TestNew_NegativeWindowReturnsError(t *testing.T) {
	_, err := New(-time.Second, nil)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNew_NilClockDefaultsToTimeNow(t *testing.T) {
	d, err := New(time.Second, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.clock == nil {
		t.Fatal("expected non-nil clock")
	}
}

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	d, _ := New(time.Minute, fixedClock(epoch))
	if !d.Allow("port:8080:opened") {
		t.Fatal("expected first call to pass")
	}
}

func TestAllow_DuplicateWithinWindowBlocked(t *testing.T) {
	d, _ := New(time.Minute, fixedClock(epoch))
	d.Allow("port:8080:opened")
	if d.Allow("port:8080:opened") {
		t.Fatal("expected duplicate within window to be blocked")
	}
}

func TestAllow_PassesAfterWindowExpires(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	d, _ := New(time.Minute, clock)
	d.Allow("port:8080:opened")
	now = epoch.Add(2 * time.Minute)
	if !d.Allow("port:8080:opened") {
		t.Fatal("expected event to pass after window expires")
	}
}

func TestAllow_ZeroWindowNeverBlocks(t *testing.T) {
	d, _ := New(0, fixedClock(epoch))
	for i := 0; i < 5; i++ {
		if !d.Allow("port:9090:opened") {
			t.Fatalf("expected zero-window deduplicator to always pass (iteration %d)", i)
		}
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	d, _ := New(time.Minute, fixedClock(epoch))
	d.Allow("port:8080:opened")
	if !d.Allow("port:9090:opened") {
		t.Fatal("expected different key to pass independently")
	}
}

func TestFlush_ResetsState(t *testing.T) {
	d, _ := New(time.Minute, fixedClock(epoch))
	d.Allow("port:8080:opened")
	d.Flush()
	if !d.Allow("port:8080:opened") {
		t.Fatal("expected key to pass after flush")
	}
}

func TestLen_TracksEntries(t *testing.T) {
	d, _ := New(time.Minute, fixedClock(epoch))
	d.Allow("a")
	d.Allow("b")
	d.Allow("a") // duplicate, should not add
	if got := d.Len(); got != 2 {
		t.Fatalf("expected Len=2, got %d", got)
	}
}
