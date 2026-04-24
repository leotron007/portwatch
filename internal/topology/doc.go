// Package topology maintains a multi-host port graph and computes
// per-host diffs between two observed states.
//
// # Graph
//
// A Graph stores the most recently observed open ports for each monitored
// host.  It is safe for concurrent use.
//
//	g := topology.New()
//	g.Set("192.168.1.1", []int{22, 80, 443})
//	ports, _ := g.Get("192.168.1.1")
//
// # Diff
//
// Diff compares two Graphs and returns per-host Change values that list
// which ports were opened and which were closed between the two snapshots.
//
//	changes := topology.Diff(prev, next)
//	for _, c := range changes {
//		fmt.Printf("%s opened=%v closed=%v\n", c.Host, c.Opened, c.Closed)
//	}
package topology
