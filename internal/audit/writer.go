package audit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileWriter opens or creates a file suitable for use as an audit log writer.
// The caller is responsible for closing the returned file.
func FileWriter(path string) (io.WriteCloser, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: mkdir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
	if err != nil {
		return nil, fmt.Errorf("audit: open file: %w", err)
	}

	return f, nil
}
