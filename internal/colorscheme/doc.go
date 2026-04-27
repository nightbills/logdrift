// Package colorscheme defines named ANSI color themes used throughout logdrift
// to style log output in the terminal.
//
// # Built-in Schemes
//
//   - default: standard 256-color palette suitable for most terminals
//   - dark:    brighter variant optimised for dark backgrounds
//   - plain:   no color codes; useful when piping output to files or other tools
//
// # Usage
//
//	scheme := colorscheme.Get(flagValue)
//	colored := scheme.Apply(scheme.Error, "something went wrong")
//
// Unknown scheme names fall back to "default" silently.
package colorscheme
