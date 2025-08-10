package main

import (
	"context"
	"os"
	"strings"

	"github.com/nduyhai/gocraft/internal/adapters/inbound/cli"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/modules/register"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/registry/embed_registry"
	"github.com/nduyhai/gocraft/internal/core/ports"
	platformlog "github.com/nduyhai/gocraft/internal/platform/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

const (
	verboseLongFlag  = "--verbose"
	verboseShortFlag = "-v"
)

// newRegistry constructs and returns the module registry with built-ins registered.
func newRegistry() ports.Registry {
	r := embed_registry.New()
	register.Builtins(r)
	return r
}

// containsVerboseFlag performs a lightweight pre-scan of args to detect verbosity.
// Cobra will do the formal parsing later; this just bootstraps logging level early.
func containsVerboseFlag(args []string) bool {
	for _, a := range args {
		if a == verboseLongFlag || a == verboseShortFlag {
			return true
		}
		// Handle combined short flags like -vv or -vfoo (presence of 'v' is enough)
		if strings.HasPrefix(a, "-") && !strings.HasPrefix(a, "--") && strings.Contains(a, "v") {
			return true
		}
	}
	return false
}

// newFxLoggerOption provides an fx.Option that configures Fx event logging
// based on the chosen verbosity and the platform logger.
func newFxLoggerOption(verbose bool) fx.Option {
	return fx.WithLogger(func() fxevent.Logger {
		if verbose {
			return &fxevent.SlogLogger{Logger: platformlog.L()}
		}
		return fxevent.NopLogger
	})
}

// runCLI builds and executes the root CLI command.
func runCLI(reg ports.Registry) error {
	root := cli.NewRootCmd(reg)
	return root.Execute()
}

func main() {
	verbose := containsVerboseFlag(os.Args[1:])
	platformlog.Init(verbose)

	app := fx.New(
		newFxLoggerOption(verbose),
		fx.Provide(
			newRegistry,
		),
		fx.Invoke(func(reg ports.Registry) error {
			return runCLI(reg)
		}),
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		os.Exit(1)
	}
	// Stop immediately after the command completes (Start returns once Invoke returns).
	_ = app.Stop(ctx)
}
