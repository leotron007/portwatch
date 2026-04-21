package sampler

import (
	"testing"
	"time"
)

const (
	minD  = 5 * time.Second
	maxD  = 60 * time.Second
	stepV = 2.0
)

func newSampler(t *testing.T) *Sampler {
	t.Helper()
	s, err := New(minD, maxD, stepV)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestNew_InvalidMin(t *testing.T) {
	_, err := New(0, maxD, stepV)
	if err == nil {
		t.Fatal("expected error for zero min")
	}
}

func TestNew_MaxLessThanMin(t *testing.T) {
	_, err := New(maxD, minD, stepV)
	if err == nil {
		t.Fatal("expected error when max < min")
	}
}

func TestNew_StepTooSmall(t *testing.T) {
	_, err := New(minD, maxD, 1.0)
	if err == nil {
		t.Fatal("expected error for step <= 1.0")
	}
}

func TestInterval_StartsAtMin(t *testing.T) {
	s := newSampler(t)
	if got := s.Interval(); got != minD {
		t.Fatalf("expected %v, got %v", minD, got)
	}
}

func TestAdvance_GrowsInterval(t *testing.T) {
	s := newSampler(t)
	s.Advance()
	if got := s.Interval(); got != 10*time.Second {
		t.Fatalf("expected 10s after one advance, got %v", got)
	}
}

func TestAdvance_CapsAtMax(t *testing.T) {
	s := newSampler(t)
	for i := 0; i < 20; i++ {
		s.Advance()
	}
	if got := s.Interval(); got != maxD {
		t.Fatalf("expected interval capped at %v, got %v", maxD, got)
	}
}

func TestRecord_ResetsToMin(t *testing.T) {
	s := newSampler(t)
	s.Advance()
	s.Advance()
	s.Record(3)
	if got := s.Interval(); got != minD {
		t.Fatalf("expected interval reset to %v after Record, got %v", minD, got)
	}
}

func TestRecord_ZeroChangesIgnored(t *testing.T) {
	s := newSampler(t)
	s.Advance() // 10s
	s.Record(0)
	if got := s.Interval(); got != 10*time.Second {
		t.Fatalf("expected interval unchanged, got %v", got)
	}
}

func TestReset_RestoresMin(t *testing.T) {
	s := newSampler(t)
	s.Advance()
	s.Advance()
	s.Reset()
	if got := s.Interval(); got != minD {
		t.Fatalf("expected %v after Reset, got %v", minD, got)
	}
}
