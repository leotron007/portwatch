// Package trendlog records periodic open-port counts to a persistent
// newline-delimited JSON file and provides lightweight trend analysis over
// configurable time windows.
//
// # Usage
//
//	tl, err := trendlog.New("/var/lib/portwatch/trend.jsonl", nil)
//	if err != nil { ... }
//
//	// Record a new observation (called each watcher tick):
//	if err := tl.Record(openPortCount, time.Now()); err != nil { ... }
//
//	// Analyse the last hour of observations:
//	summary := trendlog.Analyze(tl.Entries(), time.Hour, time.Now())
//	fmt.Println(summary.Direction) // "rising", "falling", or "stable"
package trendlog
