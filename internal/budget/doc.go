// Package budget provides a rolling-window error-budget tracker for
// portwatch.
//
// An error budget defines the maximum number of failure events that are
// tolerated within a sliding time window. Once the budget is exhausted
// the caller can take protective action such as suppressing further
// alerts, triggering an escalation, or opening a circuit breaker.
//
// Usage:
//
//	b, err := budget.New(10, time.Hour)
//	if err != nil { ... }
//
//	// on each failure:
//	b.Record()
//
//	// check health:
//	if b.Exhausted() {
//	    // take protective action
//	}
//	fmt.Println(b.Remaining()) // fraction in [0,1]
package budget
