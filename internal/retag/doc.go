// Package retag implements field renaming for structured JSON log lines.
//
// It is useful when different services emit the same semantic field under
// different names (e.g. "message" vs "msg", "svc" vs "service") and the
// user wants a normalised view across all streams.
//
// Usage:
//
//	r := retag.New(map[string]string{
//		"message": "msg",
//		"svc":     "service",
//	})
//	normalisedLine := r.Apply(rawLine)
//
// Non-JSON lines pass through unchanged. Fields absent from the log line
// are silently ignored.
package retag
