package watchdog

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

// fixedClock returns a clock whose value can be advanced manually.
func fixedClock(initial time.Time) (func() time.Time, func(time.Duration)) {
	var t = initial
	get := func() time.Time { return t }
	adv := func(d time.Duration) { t = t.Add(d) }
	return get, adv
}

func TestNew_InvalidDeadline(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero deadline")
	}
	_, err = New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative deadline")
	}
}

func TestNew_ValidDeadline(t *testing.T) {
	wd, err := New(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wd == nil {
		t.Fatal("expected non-nil watchdog")
	}
}

func TestRun_NoAlertWhenBeatsArriveInTime(t *testing.T) {
	var buf bytes.Buffer
	clock, advance := fixedClock(time.Now())

	wd, _ := New(200*time.Millisecond, WithWriter(&buf), WithClock(clock))

	beats := make(chan struct{}, 4)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go wd.Run(ctx, beats)

	// send beats before deadline expires
	for i := 0; i < 3; i++ {
		beats <- struct{}{}
		advance(50 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)
	cancel()

	if buf.Len() > 0 {
		t.Errorf("unexpected stall alert: %s", buf.String())
	}
}

func TestRun_AlertWhenStalled(t *testing.T) {
	var buf bytes.Buffer
	clock, advance := fixedClock(time.Now())

	// Use a very short deadline so the ticker fires quickly.
	wd, _ := New(100*time.Millisecond, WithWriter(&buf), WithClock(clock))

	beats := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go wd.Run(ctx, beats)

	// Advance the fake clock well past the deadline without sending beats.
	advance(500 * time.Millisecond)

	// Give the ticker goroutine time to fire at least once.
	time.Sleep(120 * time.Millisecond)
	cancel()

	if !strings.Contains(buf.String(), "stall detected") {
		t.Errorf("expected stall alert in output, got: %q", buf.String())
	}
}

func TestRun_ClosedBeatChannelExitsCleanly(t *testing.T) {
	var buf bytes.Buffer
	wd, _ := New(time.Second, WithWriter(&buf))

	beats := make(chan struct{})
	close(beats)

	ctx := context.Background()
	done := make(chan struct{})
	go func() {
		wd.Run(ctx, beats)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(time.Second):
		t.Fatal("Run did not exit after beats channel closed")
	}
}
