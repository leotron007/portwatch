package backoff

import (
	"testing"
	"time"
)

func TestNew_InvalidBase(t *testing.T) {
	_, err := New(0, time.Second, 2.0, 0)
	if err == nil {
		t.Fatal("expected error for zero base")
	}
}

func TestNew_MaxLessThanBase(t *testing.T) {
	_, err := New(time.Second, time.Millisecond, 2.0, 0)
	if err == nil {
		t.Fatal("expected error when max < base")
	}
}

func TestNew_FactorBelowOne(t *testing.T) {
	_, err := New(time.Millisecond, time.Second, 0.5, 0)
	if err == nil {
		t.Fatal("expected error for factor < 1")
	}
}

func TestNew_InvalidJitter(t *testing.T) {
	_, err := New(time.Millisecond, time.Second, 2.0, 1.5)
	if err == nil {
		t.Fatal("expected error for jitter >= 1")
	}
}

func TestNew_ValidParams(t *testing.T) {
	b, err := New(10*time.Millisecond, time.Second, 2.0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Backoff")
	}
}

func TestNext_GrowsExponentially(t *testing.T) {
	b, _ := New(10*time.Millisecond, time.Hour, 2.0, 0)

	d0 := b.Next() // attempt 0 → 10ms
	d1 := b.Next() // attempt 1 → 20ms
	d2 := b.Next() // attempt 2 → 40ms

	if d0 != 10*time.Millisecond {
		t.Errorf("d0 = %v, want 10ms", d0)
	}
	if d1 != 20*time.Millisecond {
		t.Errorf("d1 = %v, want 20ms", d1)
	}
	if d2 != 40*time.Millisecond {
		t.Errorf("d2 = %v, want 40ms", d2)
	}
}

func TestNext_CappedAtMax(t *testing.T) {
	b, _ := New(10*time.Millisecond, 15*time.Millisecond, 2.0, 0)

	b.Next() // 10ms
	d := b.Next() // would be 20ms but capped at 15ms
	if d > 15*time.Millisecond {
		t.Errorf("delay %v exceeds max 15ms", d)
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	b, _ := New(10*time.Millisecond, time.Second, 2.0, 0)
	b.Next()
	b.Next()
	if b.Attempts() != 2 {
		t.Fatalf("want 2 attempts, got %d", b.Attempts())
	}
	b.Reset()
	if b.Attempts() != 0 {
		t.Errorf("want 0 after reset, got %d", b.Attempts())
	}
	d := b.Next()
	if d != 10*time.Millisecond {
		t.Errorf("after reset d = %v, want 10ms", d)
	}
}

func TestNext_FactorOne_ConstantDelay(t *testing.T) {
	b, _ := New(5*time.Millisecond, time.Second, 1.0, 0)
	for i := 0; i < 5; i++ {
		if d := b.Next(); d != 5*time.Millisecond {
			t.Errorf("attempt %d: got %v, want 5ms", i, d)
		}
	}
}
