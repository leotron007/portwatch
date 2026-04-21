// Package jitter provides a small utility for adding randomised offsets to
// a fixed base interval.
//
// # Overview
//
// When many goroutines (or daemon ticks) fire at exactly the same wall-clock
// moment they can create short-lived resource spikes — a "thundering herd".
// Jitter spreads those wakeups across a configurable window so that the
// system load stays smooth.
//
// # Usage
//
//	j, err := jitter.New(30*time.Second, 0.2) // base=30 s, up to +6 s
//	if err != nil {
//		log.Fatal(err)
//	}
//	for {
//		time.Sleep(j.Next())
//		// … do periodic work …
//	}
package jitter
