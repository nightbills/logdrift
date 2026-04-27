// Package jsonpath implements a minimal dot-notation path extractor for JSON
// log objects.
//
// # Overview
//
// Many structured log fields are nested inside sub-objects. For example a
// Kubernetes-style log entry might look like:
//
//	{"time":"...","level":"info","meta":{"service":"api","region":"us-east-1"}}
//
// jsonpath lets CLI flags and filter expressions reference those fields with
// familiar dot notation:
//
//	--field meta.service=api
//
// # Usage
//
//	val, err := jsonpath.Extract(line, "meta.service")
//
// Only map traversal is supported. Array indexing is out of scope for the
// lightweight use-case logdrift targets.
package jsonpath
