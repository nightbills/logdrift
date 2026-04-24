// Package highlight implements keyword-based ANSI highlighting for log lines
// displayed by logdrift.
//
// Usage:
//
//	opts := highlight.Options{
//		Keywords:      []string{"error", "timeout"},
//		CaseSensitive: false,
//	}
//	formatted := highlight.Apply(line, opts)
//
// Each keyword occurrence in the line is wrapped with bold yellow ANSI escape
// sequences so it stands out in a terminal. When CaseSensitive is false
// (the default), matching is case-insensitive while preserving the original
// casing of the matched text in the output.
//
// Highlighting is applied after formatting so that ANSI codes from the
// formatter (e.g. level colours) are not disrupted.
package highlight
