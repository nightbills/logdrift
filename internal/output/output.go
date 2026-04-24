// Package output manages writing formatted log lines to one or more
// destinations (stdout, file) with optional buffering.
package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

// Writer wraps one or more io.Writers and provides thread-safe line writing.
type Writer struct {
	mu      sync.Mutex
	writers []io.Writer
	buf     *bufio.Writer
}

// New creates a Writer that writes to stdout by default.
// If filePath is non-empty, output is also tee'd to that file.
func New(filePath string) (*Writer, error) {
	writers := []io.Writer{os.Stdout}

	if filePath != "" {
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("output: open file %q: %w", filePath, err)
		}
		writers = append(writers, f)
	}

	mw := io.MultiWriter(writers...)
	return &Writer{
		writers: writers,
		buf:     bufio.NewWriter(mw),
	}, nil
}

// WriteLine writes a single formatted line followed by a newline character.
// It is safe to call from multiple goroutines.
func (w *Writer) WriteLine(line string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, err := fmt.Fprintln(w.buf, line); err != nil {
		return fmt.Errorf("output: write: %w", err)
	}
	return w.buf.Flush()
}

// Close flushes any buffered data and closes file writers if present.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.buf.Flush(); err != nil {
		return fmt.Errorf("output: flush: %w", err)
	}
	for _, wr := range w.writers {
		if c, ok := wr.(io.Closer); ok {
			if err := c.Close(); err != nil {
				return fmt.Errorf("output: close: %w", err)
			}
		}
	}
	return nil
}
