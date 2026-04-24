// Package highlight provides keyword search and ANSI highlighting
// for structured log output in logdrift.
package highlight

import (
	"strings"
)

const (
	ansiReset  = "\033[0m"
	ansiYellow = "\033[33m"
	ansiBold   = "\033[1m"
)

// Options controls how highlighting is applied.
type Options struct {
	// Keywords is the list of terms to highlight (case-insensitive).
	Keywords []string
	// CaseSensitive disables case-folding when matching keywords.
	CaseSensitive bool
}

// Apply scans line for any keyword in opts and wraps each occurrence
// with ANSI yellow+bold escape codes. Returns the (possibly modified) line.
func Apply(line string, opts Options) string {
	if len(opts.Keywords) == 0 {
		return line
	}
	for _, kw := range opts.Keywords {
		if kw == "" {
			continue
		}
		line = replaceKeyword(line, kw, opts.CaseSensitive)
	}
	return line
}

// replaceKeyword replaces all occurrences of keyword in s with a highlighted version.
func replaceKeyword(s, keyword string, caseSensitive bool) string {
	if caseSensitive {
		return strings.ReplaceAll(s, keyword, ansiBold+ansiYellow+keyword+ansiReset)
	}

	lower := strings.ToLower(s)
	lowerKW := strings.ToLower(keyword)

	var result strings.Builder
	result.Grow(len(s))

	offset := 0
	for {
		idx := strings.Index(lower[offset:], lowerKW)
		if idx == -1 {
			result.WriteString(s[offset:])
			break
		}
		abs := offset + idx
		result.WriteString(s[offset:abs])
		result.WriteString(ansiBold + ansiYellow)
		result.WriteString(s[abs : abs+len(keyword)])
		result.WriteString(ansiReset)
		offset = abs + len(keyword)
	}
	return result.String()
}
