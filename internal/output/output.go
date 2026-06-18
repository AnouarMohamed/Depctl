package output

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	mu     sync.Mutex
	quiet  bool
	out    io.Writer = os.Stdout
	errOut io.Writer = os.Stderr
)

// SetQuiet toggles quiet mode. In quiet mode, only Error output is written.
func SetQuiet(q bool) {
	mu.Lock()
	defer mu.Unlock()
	quiet = q
}

// SetWriters allows overriding the default stdout and stderr writers (useful for testing).
func SetWriters(stdout, stderr io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	out = stdout
	errOut = stderr
}

// Info prints an informational message to stdout.
func Info(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	if quiet {
		return
	}
	fmt.Fprintf(out, "info: "+format+"\n", a...)
}

// Success prints a success message to stdout.
func Success(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	if quiet {
		return
	}
	fmt.Fprintf(out, "success: "+format+"\n", a...)
}

// Warning prints a warning message to stderr.
func Warning(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	if quiet {
		return
	}
	fmt.Fprintf(errOut, "warning: "+format+"\n", a...)
}

// Error prints an error message to stderr.
func Error(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	// Error is always printed, even in quiet mode.
	fmt.Fprintf(errOut, "error: "+format+"\n", a...)
}

// Step prints a step/progress message to stdout.
func Step(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	if quiet {
		return
	}
	fmt.Fprintf(out, "--> "+format+"\n", a...)
}
