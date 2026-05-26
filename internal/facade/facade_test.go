package facade

import (
	"bytes"
	"testing"

	"github.com/goark/gocli/exitcode"
	"github.com/goark/gocli/rwi"
)

func TestExecuteUsage(t *testing.T) {
	t.Parallel()

	t.Run("no args prints usage", func(t *testing.T) {
		var out bytes.Buffer
		var errBuf bytes.Buffer

		exit := Execute(
			rwi.New(
				rwi.WithWriter(&out),
				rwi.WithErrorWriter(&errBuf),
			),
			nil,
		)

		if exit != exitcode.Normal {
			t.Fatalf("exit = %v, want %v", exit, exitcode.Normal)
		}
		if got := out.String(); got == "" || !bytes.Contains([]byte(got), []byte("Usage:")) {
			t.Fatalf("usage output missing; got %q", got)
		}
		if gotErr := errBuf.String(); gotErr != "" {
			t.Fatalf("unexpected stderr output: %q", gotErr)
		}
	})

	t.Run("subcommand help prints usage", func(t *testing.T) {
		var out bytes.Buffer
		var errBuf bytes.Buffer

		exit := Execute(
			rwi.New(
				rwi.WithWriter(&out),
				rwi.WithErrorWriter(&errBuf),
			),
			[]string{"tagslist", "-h"},
		)

		if exit != exitcode.Normal {
			t.Fatalf("exit = %v, want %v", exit, exitcode.Normal)
		}
		if got := out.String(); got == "" || !bytes.Contains([]byte(got), []byte("Usage of tagslist:")) {
			t.Fatalf("subcommand help output missing; got %q", got)
		}
		if gotErr := errBuf.String(); gotErr != "" {
			t.Fatalf("unexpected stderr output: %q", gotErr)
		}
	})
}
