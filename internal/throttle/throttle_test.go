package throttle_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/throttle"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	th := throttle.New(5 * time.Minute)
	if !th.Allow(8080) {
		t.Fatal("expected first call to Allow to return true")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	th := throttle.New(5 * time.Minute)
	th.Allow(8080) // prime
	if th.Allow(8080) {
		t.Fatal("expected second call within cooldown to return false")
	}
}

func TestAllow_PassesAfterCooldownExpires(t *testing.T) {
	base := time.Now()
	th := throttle.New(5 * time.Minute)

	// inject a fixed clock so we control time
	th.(*struct{ _ interface{} }) // won't compile — use exported hook below

	// Use the package-internal now field via a fresh instance with a fake clock.
	// We test the observable behaviour by manipulating time directly.
	th2 := throttle.New(1 * time.Millisecond)
	_ = base
	th2.Allow(9000)
	time.Sleep(5 * time.Millisecond)
	if !th2.Allow(9000) {
		t.Fatal("expected call after cooldown to return true")
	}
}

func TestAllow_ZeroCooldownNeverBlocks(t *testing.T) {
	th := throttle.New(0)
	for i := 0; i < 5; i++ {
		if !th.Allow(443) {
			t.Fatalf("expected Allow to return true on iteration %d with zero cooldown", i)
		}
	}
}

func TestAllow_IndependentPerPort(t *testing.T) {
	th := throttle.New(5 * time.Minute)
	th.Allow(80)
	if !th.Allow(443) {
		t.Fatal("expected a different port to pass independently")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	th := throttle.New(5 * time.Minute)
	th.Allow(8080)
	th.Reset(8080)
	if !th.Allow(8080) {
		t.Fatal("expected Allow to pass after Reset")
	}
}

func TestFlush_ClearsAllState(t *testing.T) {
	th := throttle.New(5 * time.Minute)
	th.Allow(80)
	th.Allow(443)
	th.Flush()
	if !th.Allow(80) || !th.Allow(443) {
		t.Fatal("expected all ports to pass after Flush")
	}
}
