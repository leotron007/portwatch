// Package circuitbreaker implements a simple circuit breaker that stops
// forwarding alerts or probe attempts after a configurable number of
// consecutive failures, giving downstream systems time to recover.
//
// States:
//
//	Closed   – normal operation; failures are counted.
//	Open     – tripped; all calls are rejected until the reset timeout expires.
//	HalfOpen – one probe call is allowed through; success closes the breaker,
//	           failure re-opens it and restarts the timeout.
package circuitbreaker

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrOpen is returned by Allow when the breaker is in the Open state.
var ErrOpen = errors.New("circuitbreaker: circuit is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // tripped; rejecting calls
	StateHalfOpen              // single probe allowed
)

// String returns a human-readable label for the state.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

// CircuitBreaker tracks consecutive failures and opens the circuit when the
// failure threshold is reached.
type CircuitBreaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	resetTimeout time.Duration
	openedAt     time.Time
	clock        func() time.Time
}

// New creates a CircuitBreaker with the given failure threshold and reset
// timeout. threshold must be ≥ 1 and resetTimeout must be > 0.
func New(threshold int, resetTimeout time.Duration) (*CircuitBreaker, error) {
	if threshold < 1 {
		return nil, fmt.Errorf("circuitbreaker: threshold must be >= 1, got %d", threshold)
	}
	if resetTimeout <= 0 {
		return nil, fmt.Errorf("circuitbreaker: resetTimeout must be > 0, got %s", resetTimeout)
	}
	return &CircuitBreaker{
		state:        StateClosed,
		threshold:    threshold,
		resetTimeout: resetTimeout,
		clock:        time.Now,
	}, nil
}

// Allow reports whether the caller may proceed. It returns ErrOpen when the
// circuit is open and the reset timeout has not yet elapsed. Once the timeout
// expires the breaker transitions to HalfOpen and allows a single probe.
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return nil
	case StateHalfOpen:
		return nil
	case StateOpen:
		if cb.clock().Sub(cb.openedAt) >= cb.resetTimeout {
			cb.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	}
	return ErrOpen
}

// RecordSuccess resets the failure counter and closes the circuit.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = StateClosed
}

// RecordFailure increments the failure counter. When the threshold is reached
// (or exceeded) the circuit transitions to Open.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= cb.threshold {
		cb.state = StateOpen
		cb.openedAt = cb.clock()
	}
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Failures returns the current consecutive failure count.
func (cb *CircuitBreaker) Failures() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.failures
}
