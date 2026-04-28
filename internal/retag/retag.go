// Package retag provides field renaming and aliasing for structured JSON log lines.
// It allows users to rename top-level JSON keys before display, for example
// mapping verbose field names to shorter canonical ones.
package retag

import (
	"encoding/json"
)

// Retagger renames fields in JSON log lines according to a mapping.
type Retagger struct {
	// mapping is oldName -> newName
	mapping map[string]string
}

// New creates a Retagger from a map of old field names to new field names.
// Keys not present in the mapping are left unchanged.
func New(mapping map[string]string) *Retagger {
	m := make(map[string]string, len(mapping))
	for k, v := range mapping {
		if k != "" && v != "" {
			m[k] = v
		}
	}
	return &Retagger{mapping: m}
}

// Apply renames fields in the given JSON line according to the mapping.
// Non-JSON lines are returned as-is. If the mapping is empty the original
// line is returned unchanged.
func (r *Retagger) Apply(line string) string {
	if len(r.mapping) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for oldKey, newKey := range r.mapping {
		val, ok := obj[oldKey]
		if !ok {
			continue
		}
		delete(obj, oldKey)
		obj[newKey] = val
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
