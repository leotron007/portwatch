// Package rotation implements size-based log file rotation for portwatch.
//
// A Rotator wraps a file path and writes to it like a standard io.Writer.
// When a write would cause the file to exceed the configured maximum size,
// the current file is renamed to <path>.1, existing backups are shifted
// (e.g. .1 → .2), and a fresh file is opened at the original path.
//
// The number of retained backup files is bounded by maxFiles; any backup
// beyond that limit is silently removed.
//
// Example:
//
//	r, err := rotation.New("/var/log/portwatch/audit.log", 10<<20, 5)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer r.Close()
//	audit.New(r) // pass as io.Writer
package rotation
