package holddown

import (
	"testing"
	"time"
)

type fixedClock struct{ t time.Time }

func (c *fixedClock) Now() time.Time { return c.t }
func (c *fixedClock) Advance(d time.Duration) { c.t = c.t.Add(d) }

func newTracker(t *testing.T, window time.Duration, c *fixedClock) *Tracker {
	t.Helper()
	tr, err := New(window, c.Now)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestNew_NegativeWindowReturnsError(t *testing.T) {
	_, err := New(-time.Second, nil)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestAllow_ZeroWindowPassesImmediately(t *testing.T) {
	c := &fixedClock{t: time.Now()}
	tr := newTracker(t, 0, c)
	if !tr.Allow("p:80", true) {
		t.Fatal("zero window should pass on first call")
	}
}

func TestAllow_FirstObservationBlocked(t *testing.T) {
	c := &fixedClock{t: time.Now()}
	tr := newTracker(t, 2*time.Second, c)
	if tr.Allow("p:443", true) {
		t.Fatal("first observation should be blocked")
	}
}

func TestAllow_PassesAfterWindowExpires(t *testing.T) {
	c := &fixedClock{t: time.Now()}
	tr := newTracker(t, 2*time.Second, c)
	tr.Allow("p:443", true)
	c.Advance(2 * time.Second)
	if !tr.Allow("p:443", true) {
		t.Fatal("should pass after window has elapsed")
	}
}

func TestAllow_AbsentResetsTimer(t *testing.T) {
	c := &fixedClock{t: time.Now()}
	tr := newTracker(t, 2*time.Second, c)
	tr.Allow("p:8080", true)
	c.Advance(3 * time.Second)
	tr.Allow("p:8080", false) // reset
	c.Advance(3 * time.Second)
	if tr.Allow("p:8080", true) {
		t.Fatal("timer should have been reset by absent call")
	}
}

func TestAllow_FalseNeverPasses(t *testing.T) {
	c := &fixedClock{t: time.Now()}
	tr := newTracker(t, time.Second, c)
	for i := 0; i < 5; i++ {
		if tr.Allow("p:22", false) {
			t.Fatal("absent key must never pass")
		}
		c.Advance(time.Second)
	}
}

func TestReset_ClearsState(t *testing.T) {
	c := &fixedClock{t: time.Now()}
	tr := newTracker(t, time.Second, c)
	tr.Allow("p:9000", true)
	c.Advance(2 * time.Second)
	tr.Reset("p:9000")
	if tr.Allow("p:9000", true) {
		t.Fatal("after Reset first observation should be blocked again")
	}
}
