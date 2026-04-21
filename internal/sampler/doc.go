// Package sampler implements adaptive scan-interval adjustment for portwatch.
//
// A Sampler starts with a minimum interval and backs off toward a configurable
// maximum when no port changes are observed, reducing unnecessary scanning
// overhead during stable periods. As soon as changes are detected the interval
// is immediately reset to the minimum so that transient activity is captured
// with high resolution.
//
// Typical usage:
//
//	s, _ := sampler.New(5*time.Second, 60*time.Second, 1.5)
//
//	for {
//		time.Sleep(s.Interval())
//		changes := scan()
//		if len(changes) > 0 {
//			s.Record(len(changes))
//		} else {
//			s.Advance()
//		}
//	}
package sampler
