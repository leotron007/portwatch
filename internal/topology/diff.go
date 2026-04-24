package topology

// Change describes a port-level change observed on a single host.
type Change struct {
	Host   string
	Opened []int
	Closed []int
}

// HasChanges returns true when there is at least one opened or closed port.
func (c Change) HasChanges() bool {
	return len(c.Opened) > 0 || len(c.Closed) > 0
}

// Diff computes per-host port changes between two Graphs.
// Hosts present only in next are treated as fully opened;
// hosts present only in prev are treated as fully closed.
func Diff(prev, next *Graph) []Change {
	prevNodes := toMap(prev.Nodes())
	nextNodes := toMap(next.Nodes())

	seen := make(map[string]struct{})
	var changes []Change

	for host, nextPorts := range nextNodes {
		seen[host] = struct{}{}
		prevPorts := prevNodes[host]
		c := Change{
			Host:   host,
			Opened: diffSlice(prevPorts, nextPorts),
			Closed: diffSlice(nextPorts, prevPorts),
		}
		if c.HasChanges() {
			changes = append(changes, c)
		}
	}

	for host, prevPorts := range prevNodes {
		if _, ok := seen[host]; ok {
			continue
		}
		changes = append(changes, Change{
			Host:   host,
			Closed: prevPorts,
		})
	}

	return changes
}

func toMap(nodes []Node) map[string][]int {
	m := make(map[string][]int, len(nodes))
	for _, n := range nodes {
		m[n.Host] = n.Ports
	}
	return m
}

// diffSlice returns elements in b that are not in a (both must be sorted).
func diffSlice(a, b []int) []int {
	set := make(map[int]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}
	var out []int
	for _, v := range b {
		if _, ok := set[v]; !ok {
			out = append(out, v)
		}
	}
	return out
}
