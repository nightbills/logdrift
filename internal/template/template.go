// Package template provides a simple line-formatting engine that renders
// structured JSON log entries using a user-supplied Go text/template string.
// Fields are extracted from the parsed JSON and exposed as a map so that
// template authors can reference arbitrary keys with {{.field}} syntax.
package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	gotemplate "text/template"
)

// Renderer applies a compiled Go template to each log line.
type Renderer struct {
	tmpl *gotemplate.Template
}

// New compiles the provided template string and returns a Renderer.
// Returns an error if the template cannot be parsed.
func New(tplStr string) (*Renderer, error) {
	tmpl, err := gotemplate.New("line").Option("missingkey=zero").Parse(tplStr)
	if err != nil {
		return nil, fmt.Errorf("template: parse: %w", err)
	}
	return &Renderer{tmpl: tmpl}, nil
}

// Render applies the template to line. If line is not valid JSON the raw line
// is returned unchanged. Template variables are top-level JSON keys; nested
// objects are represented as their JSON-encoded string form.
func (r *Renderer) Render(line string) string {
	line = strings.TrimRight(line, "\n")

	var fields map[string]any
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return line
	}

	// Flatten nested values to strings so templates stay simple.
	flat := make(map[string]string, len(fields))
	for k, v := range fields {
		switch val := v.(type) {
		case string:
			flat[k] = val
		case nil:
			flat[k] = ""
		default:
			b, _ := json.Marshal(val)
			flat[k] = string(b)
		}
	}

	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, flat); err != nil {
		return line
	}
	return buf.String()
}
