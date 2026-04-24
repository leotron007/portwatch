package graceperiod

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) Clock { return func() time.Time { return t } }

func TestNew_NegativeWindowReturnsError(t *testing.T) {
	_, err := New(-1*time.Second, nil)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNew_NilClockDefaultsToTimeNow(t *testing.T) {
	tr, err := New(time.Second, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.clock == nil {
		t.Fatal("clock should not be nil")
	}
}

func TestAllow_ZeroWindowAlwaysPasses(t *testing.T) {
	tr, _ := New(0, fixedClock(time.Now()))
	if !tr.Allow(8080) {
		t.Fatal("zero window should always allow")
	}
}

func TestAllow_BlockedBeforeWindowExpires(t *testing.T) {
	now := time.Now()
	tr, _ := New(5*time.Second, fixedClock(now))
	tr.Observe(443)
	if tr.Allow(443) {
		t.Fatal("should be blocked within grace window")
	}
}

func TestAllow_PassesAfterWindowExpires(t *testing.T) {
	now := time.Now()
	current := now
	clock := func() time.Time { return current }
	tr, _ := New(2*time.Second, clock)
	tr.Observe(443)
	current = now.Add(3 * time.Second)
	if !tr.Allow(443) {
		t.Fatal("should be allowed after grace window expires")
	}
}

func TestAllow_UnobservedPortBlocked(t *testing.T) {
	tr, _ := New(time.Second, fixedClock(time.Now()))
	if tr.Allow(9000) {
		t.Fatal("unobserved port should not be allowed")
	}
}

func TestForget_ResetsTimer(t *testing.T) {
	now := time.Now()
	current := now
	clock := func() time.Time { return current }
	tr, _ := New(2*time.Second, clock)
	tr.Observe(80)
	current = now.Add(3 * time.Second)
	if !tr.Allow(80) {
		t.Fatal("should pass after window")
	}
	tr.Forget(80)
	if tr.Allow(80) {
		t.Fatal("should be blocked again after forget")
	}
}
