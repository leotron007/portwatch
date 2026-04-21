package jitter

import (
	"testing"
	"time"
)

func TestNew_InvalidBase(t *testing.T) {
	_, err := New(0, 0.2)
	if err == nil {
		t.Fatal("expected error for zero base duration")
	}
}

func TestNew_NegativeBase(t *testing.T) {
	_, err := New(-time.Second, 0.2)
	if err == nil {
		t.Fatal("expected error for negative base duration")
	}
}

func TestNew_ZeroFraction(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero maxFraction")
	}
}

func TestNew_FractionAboveOne(t *testing.T) {
	_, err := New(time.Second, 1.5)
	if err == nil {
		t.Fatal("expected error for maxFraction > 1")
	}
}

func TestNew_ValidParams(t *testing.T) {
	j, err := New(time.Second, 0.25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j == nil {
		t.Fatal("expected non-nil Jitter")
	}
}

func TestNext_NoOffset_WhenSourceReturnsZero(t *testing.T) {
	j, _ := New(time.Second, 0.5)
	j.withSource(func() float64 { return 0 })

	got := j.Next()
	if got != time.Second {
		t.Errorf("expected %s, got %s", time.Second, got)
	}
}

func TestNext_MaxOffset_WhenSourceReturnsOne(t *testing.T) {
	base := 2 * time.Second
	fraction := 0.5
	j, _ := New(base, fraction)
	// src returns values in [0,1) but we use 0.9999… to approximate 1
	j.withSource(func() float64 { return 1.0 })

	got := j.Next()
	max := base + time.Duration(float64(base)*fraction)
	if got > max {
		t.Errorf("Next() %s exceeds expected max %s", got, max)
	}
	if got < base {
		t.Errorf("Next() %s is less than base %s", got, base)
	}
}

func TestNext_AlwaysAtLeastBase(t *testing.T) {
	base := 100 * time.Millisecond
	j, _ := New(base, 0.3)

	for i := 0; i < 50; i++ {
		if d := j.Next(); d < base {
			t.Errorf("iteration %d: Next() %s < base %s", i, d, base)
		}
	}
}

func TestNext_NeverExceedsMaximum(t *testing.T) {
	base := 100 * time.Millisecond
	fraction := 0.25
	j, _ := New(base, fraction)
	max := base + time.Duration(float64(base)*fraction)

	for i := 0; i < 50; i++ {
		if d := j.Next(); d > max {
			t.Errorf("iteration %d: Next() %s > max %s", i, d, max)
		}
	}
}
