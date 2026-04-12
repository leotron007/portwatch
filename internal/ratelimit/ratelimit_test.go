package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/ratelimit"
)

// fixedClock returns a clock whose current time can be advanced manually.
func fixedClock(initial time.Time) (*time.Time, ratelimit.Clock) {
	t := initial
	return &t, func() time.Time { return t }
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := ratelimit.New(0, 5, nil)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_InvalidMax(t *testing.T) {
	_, err := ratelimit.New(time.Second, 0, nil)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestAllow_FirstCallPasses(t *testing.T) {
	l, _ := ratelimit.New(time.Minute, 3, nil)
	if !l.Allow("port:8080") {
		t.Fatal("first call should be allowed")
	}
}

func TestAllow_BlockedAfterMax(t *testing.T) {
	now := time.Now()
	_, clock := fixedClock(now)

	l, _ := ratelimit.New(time.Minute, 2, clock)
	l.Allow("k")
	l.Allow("k")

	if l.Allow("k") {
		t.Fatal("third call within window should be blocked")
	}
}

func TestAllow_PassesAfterWindowExpires(t *testing.T) {
	now := time.Now()
	ptr, clock := fixedClock(now)

	l, _ := ratelimit.New(time.Minute, 1, clock)
	l.Allow("k")

	// Advance past the window.
	*ptr = now.Add(2 * time.Minute)

	if !l.Allow("k") {
		t.Fatal("call after window expiry should be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	now := time.Now()
	_, clock := fixedClock(now)

	l, _ := ratelimit.New(time.Minute, 1, clock)
	l.Allow("a")

	if !l.Allow("b") {
		t.Fatal("different key should not be affected")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	now := time.Now()
	_, clock := fixedClock(now)

	l, _ := ratelimit.New(time.Minute, 1, clock)
	l.Allow("k")
	l.Reset("k")

	if !l.Allow("k") {
		t.Fatal("call after Reset should be allowed")
	}
}
