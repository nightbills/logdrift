// Package redact provides field-level redaction for structured JSON log lines.
// It replaces the values of specified keys with a configurable placeholder,
// useful for scrubbing sensitive data such as tokens, passwords, or PII before
// log lines are written to output.
package redact

import (
	"encoding/json"
	"strings"
)

const defaultPlaceholder = "[REDACTED]"

// Redactor holds the set of field names to redact and the replacement string.
type Redactor struct {
	fields      map[string]struct{}
	placeholder string
}

// New returns a Redactor that will replace the value of each field in fields
// with placeholder. If placeholder is empty, "[REDACTED]" is used.
func New(fields []string, placeholder string) *Redactor {
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}
	fm := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		fm[strings.ToLower(f)] = struct{}{}
	}
	return &Redactor{fields: fm, placeholder: placeholder}
}

// Apply redacts sensitive fields from a JSON log line.
// Non-JSON lines are returned unchanged. If no fields are configured,
// the line is returned as-is.
func (r *Redactor) Apply(line string) string {
	if len(r.fields) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	modified := false
	for key := range obj {
		if _, ok := r.fields[strings.ToLower(key)]; ok {
			quoted, _ := json.Marshal(r.placeholder)
			obj[key] = json.RawMessage(quoted)
			modified = true
		}
	}

	if !modified {
		return line
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
