package stagger

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) Clock {
	return func() time.Time { return t }
}

func TestNew_InvalidCount(t *testing.T) {
	_, err := New(0, time.Second, nil)
	if err == nil {
		t.Fatal("expected error for count=0")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(2, 0, nil)
	if err == nil {
		t.Fatal("expected error for window=0")
	}
}

func TestNew_NilClockDefaultsToTimeNow(t *testing.T) {
	s, err := New(3, time.Second, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.clock == nil {
		t.Fatal("clock should not be nil")
	}
}

func TestDelay_FirstSlotIsZero(t *testing.T) {
	s, _ := New(4, 4*time.Second, fixedClock(epoch))
	if d := s.Delay(0); d != 0 {
		t.Fatalf("slot 0 delay = %v, want 0", d)
	}
}

func TestDelay_EvenlySpaced(t *testing.T) {
	s, _ := New(4, 4*time.Second, fixedClock(epoch))
	want := []time.Duration{0, time.Second, 2 * time.Second, 3 * time.Second}
	for i, w := range want {
		if got := s.Delay(i); got != w {
			t.Fatalf("slot %d delay = %v, want %v", i, got, w)
		}
	}
}

func TestDelay_SingleSlotAlwaysZero(t *testing.T) {
	s, _ := New(1, time.Minute, fixedClock(epoch))
	if d := s.Delay(0); d != 0 {
		t.Fatalf("single slot delay = %v, want 0", d)
	}
}

func TestDelay_ClampNegativeIndex(t *testing.T) {
	s, _ := New(3, 3*time.Second, fixedClock(epoch))
	if d := s.Delay(-1); d != 0 {
		t.Fatalf("clamped negative delay = %v, want 0", d)
	}
}

func TestDelay_ClampExcessIndex(t *testing.T) {
	s, _ := New(3, 3*time.Second, fixedClock(epoch))
	max := s.Delay(2)
	if d := s.Delay(99); d != max {
		t.Fatalf("clamped excess delay = %v, want %v", d, max)
	}
}

func TestSlots_LengthMatchesCount(t *testing.T) {
	s, _ := New(5, 5*time.Second, fixedClock(epoch))
	if got := len(s.Slots()); got != 5 {
		t.Fatalf("slots length = %d, want 5", got)
	}
}

func TestSlots_AnchoredToNow(t *testing.T) {
	s, _ := New(3, 3*time.Second, fixedClock(epoch))
	slots := s.Slots()
	if !slots[0].Equal(epoch) {
		t.Fatalf("first slot = %v, want %v", slots[0], epoch)
	}
	if !slots[2].Equal(epoch.Add(2 * time.Second)) {
		t.Fatalf("last slot = %v, want %v", slots[2], epoch.Add(2*time.Second))
	}
}
