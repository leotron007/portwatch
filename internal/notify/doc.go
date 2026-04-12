// Package notify provides multi-channel event dispatch for portwatch alerts.
//
// A Dispatcher fans out alert.Event values to one or more Channel
// implementations. Channels are registered at startup and called
// synchronously in registration order when Dispatch is invoked.
//
// Built-in channels:
//
//   - WriterChannel — writes a formatted line to any io.Writer (stdout, file, …).
//   - ExecChannel   — executes an external binary with event metadata
//     injected as environment variables (PORTWATCH_PORT,
//     PORTWATCH_ACTION, PORTWATCH_RULE, PORTWATCH_HOST,
//     PORTWATCH_TIME).
//
// Errors from individual channels are collected and returned as a single
// combined error so that a failing channel does not prevent others from
// receiving the event.
package notify
