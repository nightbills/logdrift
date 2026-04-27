// Package jsonpath provides a minimal dot-notation field extractor for
// structured JSON log lines. It supports simple nested key access such as
// "metadata.request.method" and returns the raw string representation of
// the value found at that path.
package jsonpath

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Extract walks a dot-separated path through a JSON object decoded from line
// and returns the value at that path as a string. If any segment of the path
// is missing, or the input is not valid JSON, an error is returned.
//
// Example:
//
//	Extract(`{"meta":{"env":"prod"}}`, "meta.env") // => "prod", nil
func Extract(line, path string) (string, error) {
	var root map[string]any
	if err := json.Unmarshal([]byte(line), &root); err != nil {
		return "", fmt.Errorf("jsonpath: invalid JSON: %w", err)
	}

	segments := strings.Split(path, ".")
	return walk(root, segments)
}

// walk recursively descends into nested maps following segments.
func walk(node map[string]any, segments []string) (string, error) {
	if len(segments) == 0 {
		return "", fmt.Errorf("jsonpath: empty path segment")
	}

	key := segments[0]
	val, ok := node[key]
	if !ok {
		return "", fmt.Errorf("jsonpath: key %q not found", key)
	}

	if len(segments) == 1 {
		return stringify(val), nil
	}

	child, ok := val.(map[string]any)
	if !ok {
		return "", fmt.Errorf("jsonpath: key %q is not an object", key)
	}

	return walk(child, segments[1:])
}

// stringify converts a JSON value to its string representation.
func stringify(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		// Avoid scientific notation for typical log field values.
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", t), "0"), ".")
	case nil:
		return ""
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
