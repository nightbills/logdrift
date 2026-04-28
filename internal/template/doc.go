// Package template provides a line-rendering engine for logdrift that formats
// structured JSON log entries using Go's text/template syntax.
//
// # Overview
//
// A Renderer is created from a template string and can be reused across many
// log lines. Each top-level JSON key is available as a template variable:
//
//	{{.level}} {{.service}} — {{.msg}}
//
// Nested objects and arrays are serialised to their compact JSON form so they
// remain accessible without requiring special helpers.
//
// Non-JSON lines pass through unchanged, making the renderer safe to use in
// pipelines that may receive mixed input.
//
// # Usage
//
//	r, err := template.New("[{{.level}}] {{.service}}: {{.msg}}")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(r.Render(line))
package template
