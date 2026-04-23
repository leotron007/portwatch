// Package retryqueue implements a bounded in-memory queue for retrying failed
// notification deliveries.
//
// # Overview
//
// When a notification channel fails to deliver an alert event, the caller may
// place the failed item into the Queue. Items are held until their NextAt
// timestamp is reached, at which point Drain returns them for re-delivery.
//
// Items that have reached the configured maximum attempt count are rejected by
// Enqueue so that the queue does not accumulate permanently undeliverable
// messages.
//
// # Usage
//
//	q, err := retryqueue.New(256, 5, time.Now)
//	if err != nil { /* handle */ }
//
//	// on failure:
//	q.Enqueue(&retryqueue.Item{
//	    ID:      eventID,
//	    Payload: encoded,
//	    NextAt:  time.Now().Add(backoff),
//	})
//
//	// in retry loop:
//	for _, item := range q.Drain() { /* re-deliver */ }
package retryqueue
