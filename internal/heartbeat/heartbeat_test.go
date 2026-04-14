package heartbeat_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/heartbeat"
)

func TestNew_InvalidInterval(t *testing.T) {
	_, err := heartbeat.New(0)
	if err == nil {
		t.Fatal("expected error for zero interval, got nil")
	}
	_, err = heartbeat.New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative interval, got nil")
	}
}

func TestNew_ValidInterval(t *testing.T) {
	tk, err := heartbeat.New(50 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tk == nil {
		t.Fatal("expected non-nil Ticker")
	}
}

func TestRun_EmitsBeats(t *testing.T) {
	tk, err := heartbeat.New(20 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	go tk.Run(ctx)

	var beats []heartbeat.Beat
	for b := range tk.C() {
		beats = append(beats, b)
	}

	if len(beats) < 2 {
		t.Fatalf("expected at least 2 beats, got %d", len(beats))
	}
}

func TestRun_SeqIsMonotonic(t *testing.T) {
	tk, err := heartbeat.New(15 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	go tk.Run(ctx)

	var prev uint64
	for b := range tk.C() {
		if b.Seq <= prev {
			t.Errorf("non-monotonic seq: got %d after %d", b.Seq, prev)
		}
		prev = b.Seq
	}
}

func TestRun_ChannelClosedAfterCancel(t *testing.T) {
	tk, err := heartbeat.New(10 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go tk.Run(ctx)
	cancel()

	// Drain and ensure channel closes within a reasonable time.
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-tk.C():
			if !ok {
				return // channel closed as expected
			}
		case <-timer.C:
			t.Fatal("channel was not closed after context cancellation")
		}
	}
}
