package state

// Diff describes ports that appeared or disappeared between two snapshots.
type Diff struct {
	Opened []int
	Closed []int
}

// IsEmpty reports whether the diff contains no changes.
func (d Diff) IsEmpty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// Compare returns the difference between a previous and a current port list.
func Compare(previous, current []int) Diff {
	prev := toSet(previous)
	curr := toSet(current)

	var d Diff
	for p := range curr {
		if !prev[p] {
			d.Opened = append(d.Opened, p)
		}
	}
	for p := range prev {
		if !curr[p] {
			d.Closed = append(d.Closed, p)
		}
	}
	return d
}

func toSet(ports []int) map[int]bool {
	m := make(map[int]bool, len(ports))
	for _, p := range ports {
		m[p] = true
	}
	return m
}
