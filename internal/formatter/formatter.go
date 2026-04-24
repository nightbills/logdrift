// Package formatter provides colorized, human-readable rendering of
// structured JSON log lines for terminal output.
package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// Options controls formatter behaviour.
type Options struct {
	NoColor    bool
	TimeFormat string // defaults to "15:04:05"
}

// Format parses a raw JSON log line and returns a pretty-printed string
// suitable for terminal display. If the line is not valid JSON it is
// returned as-is.
func Format(line string, opts Options) string {
	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return line
	}

	tf := opts.TimeFormat
	if tf == "" {
		tf = "15:04:05"
	}

	timePart := formatTime(entry, tf)
	levelPart := formatLevel(entry, opts.NoColor)
	servicePart := formatService(entry, opts.NoColor)
	msgPart := stringVal(entry, "msg", stringVal(entry, "message", ""))

	var sb strings.Builder
	if timePart != "" {
		sb.WriteString(gray(timePart, opts.NoColor))
		sb.WriteString(" ")
	}
	if servicePart != "" {
		sb.WriteString(servicePart)
		sb.WriteString(" ")
	}
	sb.WriteString(levelPart)
	if msgPart != "" {
		sb.WriteString(" ")
		sb.WriteString(msgPart)
	}
	return sb.String()
}

func formatTime(entry map[string]interface{}, layout string) string {
	raw := stringVal(entry, "time", stringVal(entry, "ts", ""))
	if raw == "" {
		return ""
	}
	for _, f := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05"} {
		if t, err := time.Parse(f, raw); err == nil {
			return t.Format(layout)
		}
	}
	return raw
}

func formatLevel(entry map[string]interface{}, noColor bool) string {
	lvl := strings.ToUpper(stringVal(entry, "level", "???"))
	if noColor {
		return fmt.Sprintf("[%s]", lvl)
	}
	var c string
	switch lvl {
	case "ERROR", "ERR", "FATAL":
		c = colorRed
	case "WARN", "WARNING":
		c = colorYellow
	case "INFO":
		c = colorCyan
	default:
		c = colorGray
	}
	return fmt.Sprintf("%s%s[%s]%s", colorBold, c, lvl, colorReset)
}

func formatService(entry map[string]interface{}, noColor bool) string {
	svc := stringVal(entry, "service", "")
	if svc == "" {
		return ""
	}
	if noColor {
		return fmt.Sprintf("(%s)", svc)
	}
	return fmt.Sprintf("%s(%s)%s", colorGray, svc, colorReset)
}

func gray(s string, noColor bool) string {
	if noColor {
		return s
	}
	return colorGray + s + colorReset
}

func stringVal(m map[string]interface{}, key, fallback string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return fallback
}
