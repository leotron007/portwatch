// Package backoff provides an exponential back-off implementation with a
// configurable base delay, ceiling, growth factor, and optional jitter.
//
// Usage:
//
//	b, err := backoff.New(100*time.Millisecond, 30*time.Second, 2.0, 0.1)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for attempt := 0; attempt < maxRetries; attempt++ {
//		if err := doWork(); err == nil {
//			b.Reset()
//			break
//		}
//		time.Sleep(b.Next())
//	}
//
// Configuration can also be loaded from YAML via backoff.FromBytes or
// backoff.FromConfig, which apply sensible defaults for any zero-valued
// fields.
package backoff
