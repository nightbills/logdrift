// Package colorscheme provides named ANSI color themes for log output.
// Themes map log levels and structural elements to ANSI escape sequences,
// allowing users to switch between built-in palettes via a CLI flag.
package colorscheme

import "fmt"

// ANSI reset sequence.
const Reset = "\033[0m"

// Scheme holds color codes for each semantic log element.
type Scheme struct {
	Name    string
	Error   string
	Warn    string
	Info    string
	Debug   string
	Key     string
	Value   string
	Time    string
	Service string
}

// Apply wraps text in the given color code and appends a reset.
func (s Scheme) Apply(color, text string) string {
	if color == "" {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// built-in schemes.
var (
	Default = Scheme{
		Name:    "default",
		Error:   "\033[31m",
		Warn:    "\033[33m",
		Info:    "\033[32m",
		Debug:   "\033[36m",
		Key:     "\033[90m",
		Value:   "\033[97m",
		Time:    "\033[90m",
		Service: "\033[35m",
	}

	Dark = Scheme{
		Name:    "dark",
		Error:   "\033[91m",
		Warn:    "\033[93m",
		Info:    "\033[92m",
		Debug:   "\033[96m",
		Key:     "\033[37m",
		Value:   "\033[97m",
		Time:    "\033[37m",
		Service: "\033[95m",
	}

	Plain = Scheme{
		Name: "plain",
	}
)

// registry maps scheme names to Scheme values.
var registry = map[string]Scheme{
	"default": Default,
	"dark":    Dark,
	"plain":   Plain,
}

// Get returns the named scheme, falling back to Default if unknown.
func Get(name string) Scheme {
	if s, ok := registry[name]; ok {
		return s
	}
	return Default
}

// Names returns all registered scheme names.
func Names() []string {
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	return out
}
