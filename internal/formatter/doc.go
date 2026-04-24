// Package formatter transforms structured JSON log entries into
// human-readable, optionally colorized lines for terminal display.
//
// # Usage
//
//	line := `{"level":"info","service":"api","msg":"ready","time":"2024-01-02T15:04:05Z"}`
//	fmt.Println(formatter.Format(line, formatter.Options{
//		NoColor:    false,
//		TimeFormat: "15:04:05",
//	}))
//
// # Output
//
// The formatter recognises the following well-known JSON keys:
//
//   - level / level  — log severity (info, warn, error, debug, fatal)
//   - msg / message  — human-readable message text
//   - service        — originating service name
//   - time / ts      — RFC 3339 timestamp
//
// Unknown keys are silently ignored; the raw line is returned unchanged
// when it cannot be parsed as JSON.
package formatter
