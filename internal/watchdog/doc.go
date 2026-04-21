// Package watchdog provides a stall detector for the portwatch scan loop.
//
// A Watchdog is initialised with a deadline duration and a writer for alerts.
// Callers pass a read-only channel of heartbeat signals (struct{}) to Run;
// every received beat resets the internal timer. If the deadline elapses
// without a beat, a human-readable stall message is written to the
// configured writer (default: os.Stderr).
//
// Typical usage:
//
//	wd, err := watchdog.New(30*time.Second, watchdog.WithWriter(logWriter))
//	if err != nil {
//		log.Fatal(err)
//	}
//	go wd.Run(ctx, heartbeatCh)
//
// The watchdog exits cleanly when ctx is cancelled or the beats channel
// is closed, making it safe to use with the daemon lifecycle.
package watchdog
