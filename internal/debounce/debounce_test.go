package debounce_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/debounce"
)

// fixedClock returns a Clock whose current time can be advanced manually.
func fixedClock(initial time.Time) (debounce.Clock, func(time.Duration)) {
	now := initial
	clock := func() time.Time { return now }
	advance := func(d time.Duration) { now = now.Add(d) }
	return clock, advance
}

func TestAllow_ZeroWindowAlwaysPasses(t *testing.T) {
	d := debounce.New(0, nil)
	if !d.Allow(8080, true) {
		t.Fatal("expected Allow to return true for zero window")
	}
}

func TestAllow_FirstObservationBlocked(t *testing.T) {
	clock, _ := fixedClock(time.Now())
	d := debounce.New(2*time.Second, clock)

	if d.Allow(8080, true) {
		t.Fatal("expected first observation to be blocked")
	}
}

func TestAllow_PassesAfterWindowExpires(t *testing.T) {
	clock, advance := fixedClock(time.Now())
	d := debounce.New(2*time.Second, clock)

	d.Allow(8080, true) // register
	advance(3 * time.Second)

	if !d.Allow(8080, true) {
		t.Fatal("expected Allow to return true after window elapsed")
	}
}

func TestAllow_StateChangeResetsTimer(t *testing.T) {
	clock, advance := fixedClock(time.Now())
	d := debounce.New(2*time.Second, clock)

	d.Allow(8080, true) // open observed
	advance(3 * time.Second)
	d.Allow(8080, false) // state flipped — timer resets

	if d.Allow(8080, false) {
		t.Fatal("expected blocked after state change reset the timer")
	}
}

func TestAllow_BlockedBeforeWindowExpires(t *testing.T) {
	clock, advance := fixedClock(time.Now())
	d := debounce.New(5*time.Second, clock)

	d.Allow(9090, true)
	advance(3 * time.Second)

	if d.Allow(9090, true) {
		t.Fatal("expected blocked: window has not yet elapsed")
	}
}

func TestReset_ClearsPendingEntry(t *testing.T) {
	clock, _ := fixedClock(time.Now())
	d := debounce.New(2*time.Second, clock)

	d.Allow(443, true)
	if d.Len() != 1 {
		t.Fatalf("expected 1 pending entry, got %d", d.Len())
	}

	d.Reset(443)
	if d.Len() != 0 {
		t.Fatalf("expected 0 pending entries after Reset, got %d", d.Len())
	}
}

func TestLen_TracksMultiplePorts(t *testing.T) {
	clock, _ := fixedClock(time.Now())
	d := debounce.New(10*time.Second, clock)

	ports := []int{80, 443, 8080}
	for _, p := range ports {
		d.Allow(p, true)
	}

	if got := d.Len(); got != len(ports) {
		t.Fatalf("expected %d pending, got %d", len(ports), got)
	}
}
