package plateau

import (
	"testing"
)

func TestNew_NegativeTolerance(t *testing.T) {
	_, err := New(-1, 3)
	if err == nil {
		t.Fatal("expected error for negative tolerance")
	}
}

func TestNew_ZeroMinRuns(t *testing.T) {
	_, err := New(0.5, 0)
	if err == nil {
		t.Fatal("expected error for zero minRuns")
	}
}

func TestNew_ValidParams(t *testing.T) {
	d, err := New(1.0, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestPlateaued_FalseBeforeEnoughRuns(t *testing.T) {
	d, _ := New(1.0, 3)
	d.Record(10.0)
	d.Record(10.5)
	if d.Plateaued() {
		t.Fatal("should not plateau after only 2 runs")
	}
}

func TestPlateaued_TrueAfterMinRuns(t *testing.T) {
	d, _ := New(1.0, 3)
	d.Record(10.0)
	d.Record(10.5)
	d.Record(9.8)
	if !d.Plateaued() {
		t.Fatal("expected plateau after 3 in-band observations")
	}
}

func TestRecord_OutOfBandResetsCounter(t *testing.T) {
	d, _ := New(1.0, 3)
	d.Record(10.0)
	d.Record(10.5)
	d.Record(20.0) // resets
	d.Record(20.3)
	if d.Plateaued() {
		t.Fatal("counter should have reset; only 2 runs since reset")
	}
}

func TestRecord_OutOfBandThenReachesMinRuns(t *testing.T) {
	d, _ := New(0.5, 2)
	d.Record(5.0)
	d.Record(50.0) // resets baseline to 50
	d.Record(50.3) // run=2, within tolerance
	if !d.Plateaued() {
		t.Fatal("expected plateau after 2 in-band runs on new baseline")
	}
}

func TestReset_ClearsState(t *testing.T) {
	d, _ := New(1.0, 2)
	d.Record(10.0)
	d.Record(10.1)
	if !d.Plateaued() {
		t.Fatal("pre-condition: should be plateaued")
	}
	d.Reset()
	if d.Plateaued() {
		t.Fatal("expected Plateaued to be false after Reset")
	}
}

func TestPlateaued_ZeroTolerance(t *testing.T) {
	d, _ := New(0, 3)
	d.Record(7.0)
	d.Record(7.0)
	d.Record(7.0)
	if !d.Plateaued() {
		t.Fatal("exact matches with zero tolerance should plateau")
	}
}
