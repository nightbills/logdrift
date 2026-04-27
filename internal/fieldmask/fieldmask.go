// Package fieldmask provides filtering of JSON log fields for display.
// It allows users to specify which top-level fields should be included
// or excluded when rendering structured log lines.
package fieldmask

import (
	"encoding/json"
	"strings"
)

// Mode controls whether the field list is an allowlist or denylist.
type Mode int

const (
	// ModeInclude renders only the specified fields.
	ModeInclude Mode = iota
	// ModeExclude renders all fields except the specified ones.
	ModeExclude
)

// Mask holds the compiled field mask configuration.
type Mask struct {
	fields map[string]struct{}
	mode   Mode
}

// New creates a Mask from a comma-separated field list and a mode.
// An empty fields string with ModeInclude means no fields are shown;
// an empty fields string with ModeExclude means all fields are shown.
func New(fields string, mode Mode) *Mask {
	m := &Mask{
		fields: make(map[string]struct{}),
		mode:   mode,
	}
	for _, f := range strings.Split(fields, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			m.fields[f] = struct{}{}
		}
	}
	return m
}

// Apply takes a raw JSON log line and returns a new JSON object containing
// only the fields permitted by the mask. Non-JSON input is returned as-is.
func (m *Mask) Apply(line string) string {
	if len(m.fields) == 0 && m.mode == ModeExclude {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	result := make(map[string]json.RawMessage, len(obj))
	for k, v := range obj {
		_, listed := m.fields[k]
		if (m.mode == ModeInclude && listed) || (m.mode == ModeExclude && !listed) {
			result[k] = v
		}
	}

	out, err := json.Marshal(result)
	if err != nil {
		return line
	}
	return string(out)
}
