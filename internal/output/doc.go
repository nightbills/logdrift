// Package output provides a thread-safe, multi-destination writer for
// logdrift's formatted log lines.
//
// # Overview
//
// A [Writer] is constructed via [New], which accepts an optional file path.
// When a file path is supplied the writer tee's every line to both stdout
// and the named file, creating or appending to it as needed.
//
// # Thread Safety
//
// [Writer.WriteLine] acquires an internal mutex before writing so it is safe
// to call concurrently from the multiple goroutines that tail individual log
// files.
//
// # Usage
//
//	w, err := output.New("/var/log/logdrift-session.log")
//	if err != nil { ... }
//	defer w.Close()
//	w.WriteLine(formattedLine)
package output
