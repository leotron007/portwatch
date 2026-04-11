package reporter

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Format controls the output format of the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Reporter writes a human-readable or structured summary of port state diffs.
type Reporter struct {
	out    io.Writer
	format Format
}

// New returns a Reporter writing to out in the given format.
// If out is nil, os.Stdout is used.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Report writes a summary of the diff to the reporter's output.
func (r *Reporter) Report(diff state.Diff) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(diff)
	default:
		return r.writeText(diff)
	}
}

func (r *Reporter) writeText(diff state.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		fmt.Fprintf(r.out, "[%s] no port changes detected\n", timestamp())
		return nil
	}
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "[%s] port changes:\n", timestamp())
	for _, p := range diff.Opened {
		fmt.Fprintf(w, "  OPENED\t%d\n", p)
	}
	for _, p := range diff.Closed {
		fmt.Fprintf(w, "  CLOSED\t%d\n", p)
	}
	return w.Flush()
}

func (r *Reporter) writeJSON(diff state.Diff) error {
	opened := portSliceJSON(diff.Opened)
	closed := portSliceJSON(diff.Closed)
	_, err := fmt.Fprintf(r.out,
		`{"timestamp":%q,"opened":[%s],"closed":[%s]}\n`,
		timestamp(), opened, closed)
	return err
}

func portSliceJSON(ports []int) string {
	if len(ports) == 0 {
		return ""
	}
	out := ""
	for i, p := range ports {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%d", p)
	}
	return out
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
