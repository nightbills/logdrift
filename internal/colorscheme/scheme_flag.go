package colorscheme

import (
	"fmt"
	"sort"
	"strings"
)

// Flag is a flag.Value-compatible type for selecting a Scheme by name.
// Use it with the standard flag package:
//
//	var sf colorscheme.Flag
//	flag.Var(&sf, "theme", sf.Usage())
type Flag struct {
	scheme Scheme
	set    bool
}

// Scheme returns the selected Scheme, defaulting to Default if not set.
func (f *Flag) Scheme() Scheme {
	if !f.set {
		return Default
	}
	return f.scheme
}

// String implements flag.Value.
func (f *Flag) String() string {
	if !f.set {
		return Default.Name
	}
	return f.scheme.Name
}

// Set implements flag.Value.
func (f *Flag) Set(value string) error {
	s, ok := registry[value]
	if !ok {
		names := Names()
		sort.Strings(names)
		return fmt.Errorf("unknown theme %q; choose from: %s", value, strings.Join(names, ", "))
	}
	f.scheme = s
	f.set = true
	return nil
}

// Usage returns a help string listing available schemes.
func (f *Flag) Usage() string {
	names := Names()
	sort.Strings(names)
	return fmt.Sprintf("color theme for log output (%s)", strings.Join(names, "|")
}
