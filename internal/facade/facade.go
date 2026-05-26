package facade

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/goark/errs"
	"github.com/goark/gocli/exitcode"
	"github.com/goark/gocli/rwi"
	"github.com/goark/struct2pflag"
	"github.com/spf13/pflag"

	"github.com/spiegel-im-spiegel/tagtools/internal/tagslist"
	"github.com/spiegel-im-spiegel/tagtools/internal/toptags"
	"github.com/spiegel-im-spiegel/tagtools/internal/verify"
)

var (
	// Name is application name.
	Name = "tagtools"
)

type jsonOutputError struct {
	Payload string
}

func (e *jsonOutputError) Error() string {
	return "verification failed"
}

// Execute is called from main function.
func Execute(ui *rwi.RWI, args []string) (exit exitcode.ExitCode) {
	defer func() {
		if r := recover(); r != nil {
			_ = ui.OutputErrln("Panic:", r)
			for depth := 0; ; depth++ {
				pc, _, line, ok := runtime.Caller(depth)
				if !ok {
					break
				}
				_ = ui.OutputErrln(" ->", depth, ":", runtime.FuncForPC(pc).Name(), ": line", line)
			}
			exit = exitcode.Abnormal
		}
	}()

	exit = exitcode.Normal
	if err := run(ui, args); err != nil {
		var jerr *jsonOutputError
		if errors.As(err, &jerr) {
			_ = ui.OutputErrln(jerr.Payload)
			exit = exitcode.Abnormal
			return
		}
		_ = ui.OutputErrln(fmt.Sprintf("error: %v", err))
		exit = exitcode.Abnormal
	}
	return
}

func run(ui *rwi.RWI, args []string) error {
	if len(args) == 0 {
		printUsage(ui)
		return nil
	}

	subcmd := args[0]
	subArgs := args[1:]

	switch subcmd {
	case "tagslist":
		cfg := tagslist.DefaultConfig()
		fs := pflag.NewFlagSet("tagslist", pflag.ContinueOnError)
		fs.SetOutput(ui.Writer())
		struct2pflag.Bind(fs, &cfg)
		if err := fs.Parse(subArgs); err != nil {
			if errs.Is(err, pflag.ErrHelp) {
				return nil
			}
			return errs.Wrap(err)
		}
		if err := tagslist.Run(cfg); err != nil {
			return errs.Wrap(err)
		}
		_ = ui.Outputln("Updated", cfg.Out)
		return nil
	case "toptags":
		cfg := toptags.DefaultConfig()
		fs := pflag.NewFlagSet("toptags", pflag.ContinueOnError)
		fs.SetOutput(ui.Writer())
		struct2pflag.Bind(fs, &cfg)
		if err := fs.Parse(subArgs); err != nil {
			if errs.Is(err, pflag.ErrHelp) {
				return nil
			}
			return errs.Wrap(err)
		}
		if err := toptags.Run(cfg); err != nil {
			return errs.Wrap(err)
		}
		_ = ui.Outputln("Updated", cfg.Out)
		return nil
	case "all":
		if err := toptags.Run(toptags.DefaultConfig()); err != nil {
			return errs.Wrap(err)
		}
		if err := tagslist.Run(tagslist.DefaultConfig()); err != nil {
			return errs.Wrap(err)
		}
		_ = ui.Outputln("Updated data/toptags.json and .github/workflows/tagslist.csv")
		return nil
	case "verify":
		cfg := verify.DefaultConfig()
		fs := pflag.NewFlagSet("verify", pflag.ContinueOnError)
		fs.SetOutput(ui.Writer())
		struct2pflag.Bind(fs, &cfg)
		if err := fs.Parse(subArgs); err != nil {
			if errs.Is(err, pflag.ErrHelp) {
				return nil
			}
			return errs.Wrap(err)
		}
		if err := verify.Run(cfg); err != nil {
			if cfg.Debug {
				payload, jerr := verify.JSONError(err, cfg)
				if jerr != nil {
					return errs.Wrap(jerr)
				}
				return &jsonOutputError{Payload: payload}
			}
			return errs.Wrap(err)
		}
		if cfg.Debug {
			payload, jerr := verify.JSONSuccess(cfg)
			if jerr != nil {
				return errs.Wrap(jerr)
			}
			_ = ui.Outputln(payload)
			return nil
		}
		_ = ui.Outputln("Verification passed")
		return nil
	case "help", "-h", "--help":
		printUsage(ui)
		return nil
	default:
		return errs.New("unknown subcommand", errs.WithContext("subcommand", subcmd))
	}
}

func printUsage(ui *rwi.RWI) {
	_ = ui.Outputln("tagtools: helper CLI for tagslist/toptags")
	_ = ui.Outputln("")
	_ = ui.Outputln("Usage:")
	_ = ui.Outputln("  tagtools tagslist [flags]")
	_ = ui.Outputln("  tagtools toptags [flags]")
	_ = ui.Outputln("  tagtools all")
	_ = ui.Outputln("  tagtools verify [flags]")
}
