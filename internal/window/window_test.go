package window_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/window"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) window.Clock {
	return func() time.Time { return t }
}

func TestNew_InvalidSize(t *testing.T) {
	_, err := window.New(0, 4, nil)
	if err == nil {
		t.Fatal("expected error for zero size")
	}
}

func TestNew_InvalidBuckets(t *testing.T) {
	_, err := window.New(time.Minute, 0, nil)
	if err == nil {
		t.Fatal("expected error for zero buckets")
	}
}

func TestNew_ValidParams(t *testing.T) {
	w, err := window.New(time.Minute, 6, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestCount_EmptyWindowIsZero(t *testing.T) {
	w, _ := window.New(time.Minute, 4, fixedClock(epoch))
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_IncrementsCount(t *testing.T) {
	w, _ := window.New(time.Minute, 4, fixedClock(epoch))
	w.Add(3)
	w.Add(2)
	if got := w.Count(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestCount_ExpiredBucketsNotCounted(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	w, _ := window.New(time.Minute, 4, clock)
	w.Add(10)
	// Advance beyond the window size.
	now = epoch.Add(2 * time.Minute)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after window expired, got %d", got)
	}
}

func TestReset_ClearsAllCounts(t *testing.T) {
	w, _ := window.New(time.Minute, 4, fixedClock(epoch))
	w.Add(7)
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCount_MultipleBuckets(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	w, _ := window.New(time.Minute, 4, clock) // 4 buckets × 15 s each

	w.Add(1)                       // bucket at t=0
	now = epoch.Add(20 * time.Second)
	w.Add(2)                       // bucket at t=20s
	now = epoch.Add(40 * time.Second)
	w.Add(3)                       // bucket at t=40s

	// All three buckets are within 60 s of t=40s.
	if got := w.Count(); got != 6 {
		t.Fatalf("expected 6, got %d", got)
	}
}
