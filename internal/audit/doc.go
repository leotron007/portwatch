// Package audit provides a structured JSON audit logger for portwatch.
//
// Each to Logger.Log appends a newline-delimited JSON entry to the
// configured writer containing a UTC timestamp, event name,
// protocol, and an optional human-readable reason.
//
// Typical usage:
//
//	logger := audit.New(w)
//	logger.Log("port_opened", 8080, "tcp", "new listener detected")
//
// Use FileWriter to obtain an append-only file handle:
//
//	w, err := audit.FileWriter("/var/log/portwatch/audit.log")
//	if err != nil { ... }
//	defer w.Close()
package audit
