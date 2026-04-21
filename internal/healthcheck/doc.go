// Package healthcheck provides a lightweight health-reporting mechanism for
// the portwatch daemon.
//
// A [Reporter] holds a set of named [Check] probes. Each probe is a function
// that returns a non-nil error when the component it monitors is unhealthy.
// Checks may be marked Critical; a single critical failure raises the overall
// status to [StatusUnhealthy], while non-critical failures result in
// [StatusDegraded].
//
// Usage:
//
//	r := healthcheck.New(os.Stdout)
//	r.Register(healthcheck.Check{
//		Name:     "state-file",
//		Critical: true,
//		Probe:    func() error { return stateStore.Ping() },
//	})
//	status := r.Write() // prints summary and returns overall Status
package healthcheck
