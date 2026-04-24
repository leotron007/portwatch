package flap_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/flap"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := flap.New(0, 3, nil)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_ThresholdBelowTwo(t *testing.T) {
	_, err := flap.New(time.Minute, 1, nil)
	if err == nil {
		t.Fatal("expected error for threshold < 2")
	}
}

func TestNew_ValidParams(t *testing.T) {
	d, err := flap.New(time.Minute, 3, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Detector")
	}
}

func TestRecord_BelowThresholdNotFlapping(t *testing.T) {
	now := epoch
	d, _ := flap.New(time.Minute, 3, fixedClock(now))

	if d.Record(8080) {
		t.Error("single event should not be flapping")
	}
	if d.Record(8080) {
		t.Error("two events below threshold should not be flapping")
	}
}

func TestRecord_AtThresholdIsFlapping(t *testing.T) {
	now := epoch
	d, _ := flap.New(time.Minute, 3, fixedClock(now))

	d.Record(8080)
	d.Record(8080)
	if !d.Record(8080) {
		t.Error("third event should trigger flapping")
	}
}

func TestRecord_OldEventsExpire(t *testing.T) {
	now := epoch
	clock := &now
	d, _ := flap.New(time.Minute, 3, func() time.Time { return *clock })

	d.Record(9000)
	d.Record(9000)
	// Advance past the window so previous events expire.
	*clock = now.Add(2 * time.Minute)
	if d.Record(9000) {
		t.Error("expired events should not count toward threshold")
	}
}

func TestReset_ClearsHistory(t *testing.T) {
	now := epoch
	d, _ := flap.New(time.Minute, 2, fixedClock(now))

	d.Record(443)
	d.Reset(443)
	if d.Record(443) {
		t.Error("after reset a single event should not be flapping")
	}
}

func TestFlapping_ReturnsFlapPorts(t *testing.T) {
	now := epoch
	d, _ := flap.New(time.Minute, 2, fixedClock(now))

	d.Record(80)
	d.Record(80)
	d.Record(443) // only one event

	ports := d.Flapping()
	if len(ports) != 1 || ports[0] != 80 {
		t.Errorf("expected [80], got %v", ports)
	}
}
