package budget

import (
	"testing"
	"time"
)

type fixedClock struct{ t time.Time }

func (f *fixedClock) Now() time.Time { return f.t }

func newBudget(t *testing.T, capacity int, window time.Duration) *Budget {
	t.Helper()
	b, err := New(capacity, window)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return b
}

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := New(0, time.Minute)
	if err == nil {
		t.Fatal("expected error for capacity=0")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(5, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestRemaining_StartsAtOne(t *testing.T) {
	b := newBudget(t, 4, time.Minute)
	if got := b.Remaining(); got != 1.0 {
		t.Fatalf("want 1.0 got %v", got)
	}
}

func TestRecord_DecreasesRemaining(t *testing.T) {
	b := newBudget(t, 4, time.Minute)
	b.Record()
	b.Record()
	want := 0.5
	if got := b.Remaining(); got != want {
		t.Fatalf("want %v got %v", want, got)
	}
}

func TestExhausted_WhenCapacityReached(t *testing.T) {
	b := newBudget(t, 2, time.Minute)
	b.Record()
	b.Record()
	if !b.Exhausted() {
		t.Fatal("expected budget to be exhausted")
	}
}

func TestRecord_PrunesExpiredEvents(t *testing.T) {
	b := newBudget(t, 3, time.Minute)
	clk := &fixedClock{t: time.Now()}
	b.clock = clk.Now

	b.Record()
	b.Record()

	// advance past window
	clk.t = clk.t.Add(2 * time.Minute)
	b.Record() // triggers prune; only 1 event remains

	if b.Exhausted() {
		t.Fatal("budget should not be exhausted after prune")
	}
	want := float64(2) / float64(3)
	if got := b.Remaining(); got != want {
		t.Fatalf("want %v got %v", want, got)
	}
}

func TestFromBytes_Valid(t *testing.T) {
	yaml := []byte("capacity: 5\nwindow: 30m\n")
	b, err := FromBytes(yaml)
	if err != nil {
		t.Fatalf("FromBytes: %v", err)
	}
	if b.capacity != 5 {
		t.Fatalf("want capacity 5 got %d", b.capacity)
	}
}

func TestFromBytes_InvalidWindow(t *testing.T) {
	yaml := []byte("capacity: 5\nwindow: notaduration\n")
	_, err := FromBytes(yaml)
	if err == nil {
		t.Fatal("expected error for invalid window")
	}
}

func TestString_ContainsSummary(t *testing.T) {
	b := newBudget(t, 10, time.Hour)
	b.Record()
	s := b.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
