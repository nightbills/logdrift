// Package truncate provides utilities for truncating long log lines
// to a configurable maximum width, preserving ANSI escape codes.
package truncate

import (
	"strings"
	"unicode/utf8"
)

const (
	// DefaultMaxWidth is the default maximum number of visible characters per line.
	DefaultMaxWidth = 200
	// ellipsis is appended when a line is truncated.
	ellipsis = "…"
)

// Truncator holds configuration for line truncation.
type Truncator struct {
	maxWidth int
}

// New returns a Truncator that clips visible characters to maxWidth.
// If maxWidth is <= 0, truncation is disabled and lines are returned as-is.
func New(maxWidth int) *Truncator {
	return &Truncator{maxWidth: maxWidth}
}

// Line truncates s to at most t.maxWidth visible (non-ANSI) characters,
// appending an ellipsis when truncation occurs. ANSI escape sequences
// are counted as zero-width so colour output is preserved up to the cut.
func (t *Truncator) Line(s string) string {
	if t.maxWidth <= 0 {
		return s
	}

	visible := 0
	i := 0
	for i < len(s) {
		// Detect ANSI escape sequence: ESC '[' ... letter
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			end := strings.IndexAny(s[i+2:], "ABCDEFGHJKSTfmnsulh")
			if end >= 0 {
				i += 2 + end + 1
				continue
			}
		}

		_, size := utf8.DecodeRuneInString(s[i:])
		visible++
		if visible > t.maxWidth {
			return s[:i] + ellipsis
		}
		i += size
	}
	return s
}
