// Package flap provides a flap detector for monitored ports.
//
// A port is considered "flapping" when it transitions between open and
// closed states too frequently within a configurable time window.  This
// can indicate an unstable service, a misconfigured firewall rule, or a
// port scan in progress.
//
// Basic usage:
//
//	detector, err := flap.New(30*time.Second, 4, nil)
//	if err != nil { ... }
//
//	// Call Record each time a port's state changes.
//	if detector.Record(port) {
//		log.Printf("port %d is flapping", port)
//	}
//
// The zero clock argument defaults to time.Now.
package flap
