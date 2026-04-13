package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Digest represents a deterministic hash of a set of open ports.
type Digest struct {
	hash  string
	ports []int
}

// New computes a SHA-256 digest from the given list of open ports.
// The port list is sorted before hashing to ensure determinism.
func New(ports []int) *Digest {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	parts := make([]string, len(sorted))
	for i, p := range sorted {
		parts[i] = fmt.Sprintf("%d", p)
	}

	input := strings.Join(parts, ",")
	sum := sha256.Sum256([]byte(input))

	return &Digest{
		hash:  hex.EncodeToString(sum[:]),
		ports: sorted,
	}
}

// Hash returns the hex-encoded SHA-256 hash string.
func (d *Digest) Hash() string {
	return d.hash
}

// Ports returns the sorted port list used to compute the digest.
func (d *Digest) Ports() []int {
	out := make([]int, len(d.ports))
	copy(out, d.ports)
	return out
}

// Equal reports whether two digests represent the same port set.
func (d *Digest) Equal(other *Digest) bool {
	if other == nil {
		return false
	}
	return d.hash == other.hash
}

// String returns a short prefix of the hash suitable for display.
func (d *Digest) String() string {
	if len(d.hash) >= 12 {
		return d.hash[:12]
	}
	return d.hash
}
