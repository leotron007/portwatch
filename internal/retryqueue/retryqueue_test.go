package retryqueue

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_Defaults(t *testing.T) {
	q, err := New(0, 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.maxSize != 256 || q.maxTries != 5 {
		t.Fatalf("expected defaults 256/5, got %d/%d", q.maxSize, q.maxTries)
	}
}

func TestEnqueue_AddsItem(t *testing.T) {
	q, _ := New(10, 3, fixedClock(epoch))
	item := &Item{ID: "a", Payload: []byte("hello")}
	if err := q.Enqueue(item); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Len() != 1 {
		t.Fatalf("expected len 1, got %d", q.Len())
	}
}

func TestEnqueue_FullQueueReturnsError(t *testing.T) {
	q, _ := New(2, 5, fixedClock(epoch))
	for i := 0; i < 2; i++ {
		_ = q.Enqueue(&Item{ID: "x"})
	}
	if err := q.Enqueue(&Item{ID: "overflow"}); err == nil {
		t.Fatal("expected error for full queue")
	}
}

func TestEnqueue_ExceededAttemptsReturnsError(t *testing.T) {
	q, _ := New(10, 3, fixedClock(epoch))
	item := &Item{ID: "b", Attempts: 3}
	if err := q.Enqueue(item); err == nil {
		t.Fatal("expected error when attempts >= maxTries")
	}
}

func TestEnqueue_NilItemReturnsError(t *testing.T) {
	q, _ := New(10, 3, fixedClock(epoch))
	if err := q.Enqueue(nil); err == nil {
		t.Fatal("expected error for nil item")
	}
}

func TestDrain_ReturnsOnlyReadyItems(t *testing.T) {
	now := epoch
	q, _ := New(10, 5, fixedClock(now))

	_ = q.Enqueue(&Item{ID: "ready", NextAt: now.Add(-time.Second)})
	_ = q.Enqueue(&Item{ID: "future", NextAt: now.Add(time.Minute)})

	ready := q.Drain()
	if len(ready) != 1 || ready[0].ID != "ready" {
		t.Fatalf("expected 1 ready item, got %d", len(ready))
	}
	if q.Len() != 1 {
		t.Fatalf("expected 1 remaining item, got %d", q.Len())
	}
}

func TestDrain_EmptyQueueReturnsNil(t *testing.T) {
	q, _ := New(10, 5, fixedClock(epoch))
	if got := q.Drain(); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestDrain_RemovesReturnedItems(t *testing.T) {
	now := epoch
	q, _ := New(10, 5, fixedClock(now))
	_ = q.Enqueue(&Item{ID: "c", NextAt: now})
	q.Drain()
	if q.Len() != 0 {
		t.Fatalf("expected empty queue after drain, got %d", q.Len())
	}
}
