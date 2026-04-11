package filter

// Filter holds criteria for including or excluding ports from scan results.
type Filter struct {
	// AllowedPorts is an explicit list of ports to include; empty means all.
	AllowedPorts []int
	// IgnoredPorts is a list of ports to always exclude from results.
	IgnoredPorts []int
	// MinPort is the lower bound of the port range (inclusive).
	MinPort int
	// MaxPort is the upper bound of the port range (inclusive).
	MaxPort int
}

// New returns a Filter with sensible defaults (all ports 1–65535, nothing ignored).
func New() *Filter {
	return &Filter{
		MinPort: 1,
		MaxPort: 65535,
	}
}

// Allow returns true when port p should be included in results.
func (f *Filter) Allow(p int) bool {
	if p < f.MinPort || p > f.MaxPort {
		return false
	}
	for _, ignored := range f.IgnoredPorts {
		if ignored == p {
			return false
		}
	}
	if len(f.AllowedPorts) == 0 {
		return true
	}
	for _, allowed := range f.AllowedPorts {
		if allowed == p {
			return true
		}
	}
	return false
}

// Apply filters a slice of ports, returning only those that pass Allow.
func (f *Filter) Apply(ports []int) []int {
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if f.Allow(p) {
			out = append(out, p)
		}
	}
	return out
}
