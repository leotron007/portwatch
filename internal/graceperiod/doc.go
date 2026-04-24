// Package graceperiod provides a Tracker that suppresses port-change alerts
// until a port has been continuously observed for a configurable window.
//
// This avoids alert storms caused by ephemeral ports (e.g. short-lived
// outbound connections) that open and close within a few seconds.
//
// Usage:
//
//	tr, err := graceperiod.New(3*time.Second, nil)
//	tr.Observe(port)
//	if tr.Allow(port) {
//	    // emit alert
//	}
//
// When a port is no longer present, call Forget so the timer resets if the
// port reappears later.
package graceperiod
