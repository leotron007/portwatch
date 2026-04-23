// Package window implements a sliding-window counter that partitions a
// configurable time span into discrete buckets.  Each bucket accumulates
// event counts for its sub-interval; buckets whose age exceeds the total
// window size are treated as expired and excluded from Count.
//
// Typical usage:
//
//	w, err := window.New(time.Minute, 6, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	w.Add(1)           // record an event
//	total := w.Count() // events in the last minute
//
// The zero value is not usable; always construct via New.
package window
