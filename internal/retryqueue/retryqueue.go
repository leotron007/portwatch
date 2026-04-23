// Package retryqueue provides a bounded, persistent retry queue for failed
// alert notifications. Items that fail delivery are re-enqueued with
// exponential back-off until they succeed or exceed the maximum attempt count.
package retryqueue

import (
	"errors"
	"sync"
	"time"
)

// Item holds a single queued notification payload together with its delivery
// metadata.
type Item struct {
	ID       string
	Payload  []byte
	Attempts int
	NextAt   time.Time
}

// Queue is a thread-safe, in-memory retry queue.
type Queue struct {
	mu       sync.Mutex
	items    []*Item
	maxSize  int
	maxTries int
	clock    func() time.Time
}

// New creates a Queue with the given capacity and maximum attempt limits.
// maxSize <= 0 defaults to 256; maxTries <= 0 defaults to 5.
func New(maxSize, maxTries int, clock func() time.Time) (*Queue, error) {
	if maxSize <= 0 {
		maxSize = 256
	}
	if maxTries <= 0 {
		maxTries = 5
	}
	if clock == nil {
		clock = time.Now
	}
	return &Queue{maxSize: maxSize, maxTries: maxTries, clock: clock}, nil
}

// Enqueue adds item to the queue. It returns an error when the queue is full
// or the item has already exceeded the maximum attempt count.
func (q *Queue) Enqueue(item *Item) error {
	if item == nil {
		return errors.New("retryqueue: nil item")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if item.Attempts >= q.maxTries {
		return errors.New("retryqueue: item exceeded max attempts")
	}
	if len(q.items) >= q.maxSize {
		return errors.New("retryqueue: queue full")
	}
	q.items = append(q.items, item)
	return nil
}

// Drain returns all items whose NextAt time is at or before now and removes
// them from the queue.
func (q *Queue) Drain() []*Item {
	now := q.clock()
	q.mu.Lock()
	defer q.mu.Unlock()
	var ready, remaining []*Item
	for _, it := range q.items {
		if !it.NextAt.After(now) {
			ready = append(ready, it)
		} else {
			remaining = append(remaining, it)
		}
	}
	q.items = remaining
	return ready
}

// Len returns the current number of queued items.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}
