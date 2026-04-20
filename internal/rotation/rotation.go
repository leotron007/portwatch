// Package rotation provides log file rotation for portwatch output files.
// It wraps an underlying file writer and rotates when the file exceeds a
// configured maximum size, keeping a bounded number of backup files.
package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Rotator writes to a file and rotates it when it exceeds MaxBytes.
type Rotator struct {
	mu       sync.Mutex
	path     string
	maxBytes int64
	maxFiles int
	file     *os.File
	size     int64
}

// New creates a Rotator that writes to path. maxBytes is the size threshold
// that triggers rotation; maxFiles is the number of backup files to keep.
// Returns an error if maxBytes or maxFiles is not positive.
func New(path string, maxBytes int64, maxFiles int) (*Rotator, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("rotation: maxBytes must be positive, got %d", maxBytes)
	}
	if maxFiles <= 0 {
		return nil, fmt.Errorf("rotation: maxFiles must be positive, got %d", maxFiles)
	}
	r := &Rotator{path: path, maxBytes: maxBytes, maxFiles: maxFiles}
	if err := r.openOrCreate(); err != nil {
		return nil, err
	}
	return r, nil
}

// Write implements io.Writer. It rotates the file if writing p would exceed
// the configured size limit.
func (r *Rotator) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.size+int64(len(p)) > r.maxBytes {
		if err := r.rotate(); err != nil {
			return 0, fmt.Errorf("rotation: rotate: %w", err)
		}
	}
	n, err := r.file.Write(p)
	r.size += int64(n)
	return n, err
}

// Close closes the underlying file.
func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.file == nil {
		return nil
	}
	err := r.file.Close()
	r.file = nil
	return err
}

func (r *Rotator) openOrCreate() error {
	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("rotation: open %s: %w", r.path, err)
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return fmt.Errorf("rotation: stat %s: %w", r.path, err)
	}
	r.file = f
	r.size = info.Size()
	return nil
}

func (r *Rotator) rotate() error {
	if err := r.file.Close(); err != nil {
		return err
	}
	r.file = nil

	// Shift existing backups: .2 -> .3, .1 -> .2, etc.
	for i := r.maxFiles - 1; i >= 1; i-- {
		old := fmt.Sprintf("%s.%d", r.path, i)
		new := fmt.Sprintf("%s.%d", r.path, i+1)
		_ = os.Rename(old, new)
	}
	// Remove the oldest if it overflows.
	_ = os.Remove(fmt.Sprintf("%s.%d", r.path, r.maxFiles+1))

	if err := os.Rename(r.path, filepath.FromSlash(r.path+".1")); err != nil {
		return fmt.Errorf("rotation: rename current log: %w", err)
	}
	return r.openOrCreate()
}
