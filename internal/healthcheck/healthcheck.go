package healthcheck

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Status represents the overall health of the daemon.
type Status int

const (
	StatusOK      Status = iota
	StatusDegraded        // some checks failed but daemon is running
	StatusUnhealthy       // critical checks failed
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusDegraded:
		return "degraded"
	default:
		return "unhealthy"
	}
}

// Check is a named health probe that returns an error when unhealthy.
type Check struct {
	Name     string
	Critical bool
	Probe    func() error
}

// Result holds the outcome of a single check.
type Result struct {
	Name     string    `json:"name"`
	Critical bool      `json:"critical"`
	Healthy  bool      `json:"healthy"`
	Message  string    `json:"message,omitempty"`
	CheckedAt time.Time `json:"checked_at"`
}

// Reporter runs registered checks and writes a summary.
type Reporter struct {
	checks []Check
	w      io.Writer
}

// New returns a Reporter that writes to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Register adds a check to the reporter.
func (r *Reporter) Register(c Check) {
	r.checks = append(r.checks, c)
}

// Run executes all checks and returns the results and overall status.
func (r *Reporter) Run() ([]Result, Status) {
	now := time.Now().UTC()
	results := make([]Result, 0, len(r.checks))
	overall := StatusOK

	for _, c := range r.checks {
		res := Result{Name: c.Name, Critical: c.Critical, Healthy: true, CheckedAt: now}
		if err := c.Probe(); err != nil {
			res.Healthy = false
			res.Message = err.Error()
			if c.Critical {
				overall = StatusUnhealthy
			} else if overall != StatusUnhealthy {
				overall = StatusDegraded
			}
		}
		results = append(results, res)
	}
	return results, overall
}

// Write runs all checks and writes a human-readable summary to the reporter's writer.
func (r *Reporter) Write() Status {
	results, status := r.Run()
	fmt.Fprintf(r.w, "health:n	for _, res := range results {
		mark := "✓"
		if !res.Healthy {
			mark = "✗"
		}
		line := fmt.Sprintf("  [%s] %s", mark, res.Name)
		if res.Message != "" {
			line += ": " + res.Message
		}
		fmt.Fprintln(r.w, line)
	}
	return status
}
