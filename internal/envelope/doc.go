// Package envelope defines the Envelope type, which wraps a port-change
// event with routing metadata such as severity, host, timestamp, and
// arbitrary string labels.
//
// Envelopes are the unit of communication between the watcher layer and
// downstream consumers such as notifiers, auditors, and reporters.  By
// centralising metadata in a single struct the rest of the codebase can
// handle events uniformly without coupling to raw port integers.
//
// Typical usage:
//
//	e := envelope.New(host, port, envelope.EventOpened, envelope.SeverityWarning).
//		WithLabel("env", "production")
//	 dispatcher.Dispatch(ctx, e)
package envelope
