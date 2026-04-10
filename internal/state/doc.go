// Package state manages the persistence and comparison of port scan snapshots.
//
// A Snapshot records which TCP ports were open at a specific moment in time.
// The Store type serialises snapshots to a JSON file so that portwatch can
// detect changes across process restarts.
//
// Usage:
//
//	store, err := state.New("/var/lib/portwatch/state.json")
//	if err != nil { /* handle */ }
//
//	prev := store.Current()
//	// ... run scanner ...
//	diff := state.Compare(prev.Ports, newPorts)
//	if !diff.IsEmpty() { /* alert */ }
//	store.Save(newPorts)
package state
