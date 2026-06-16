package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestOutputFunctions(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	SetWriters(stdout, stderr)
	SetQuiet(false)

	t.Run("Info", func(t *testing.T) {
		stdout.Reset()
		Info("hello %s", "world")
		got := stdout.String()
		want := "info: hello world\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Success", func(t *testing.T) {
		stdout.Reset()
		Success("hello %d", 123)
		got := stdout.String()
		want := "success: hello 123\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Step", func(t *testing.T) {
		stdout.Reset()
		Step("next phase")
		got := stdout.String()
		want := "--> next phase\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Warning", func(t *testing.T) {
		stderr.Reset()
		Warning("be careful")
		got := stderr.String()
		want := "warning: be careful\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Error", func(t *testing.T) {
		stderr.Reset()
		Error("something went wrong")
		got := stderr.String()
		want := "error: something went wrong\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestQuietMode(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	SetWriters(stdout, stderr)
	SetQuiet(true)

	t.Run("Suppress Info, Success, Step, Warning", func(t *testing.T) {
		stdout.Reset()
		stderr.Reset()

		Info("info message")
		Success("success message")
		Step("step message")
		Warning("warning message")

		if stdout.Len() > 0 {
			t.Errorf("unexpected output in stdout: %q", stdout.String())
		}
		if stderr.Len() > 0 {
			t.Errorf("unexpected output in stderr: %q", stderr.String())
		}
	})

	t.Run("Allow Error even if quiet", func(t *testing.T) {
		stderr.Reset()
		Error("error message")
		got := stderr.String()
		if !strings.Contains(got, "error: error message") {
			t.Errorf("expected error to be printed even in quiet mode, got %q", got)
		}
	})
}
