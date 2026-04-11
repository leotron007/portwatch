// Package reporter provides formatted output of port-state diff results.
//
// It supports two output formats:
//
//   - FormatText: a human-readable tabular summary suitable for terminal output.
//   - FormatJSON: a single-line JSON object suitable for log ingestion pipelines.
//
// Example usage:
//
//	r := reporter.New(os.Stdout, reporter.FormatText)
//	if err := r.Report(diff); err != nil {
//		log.Fatal(err)
//	}
package reporter
