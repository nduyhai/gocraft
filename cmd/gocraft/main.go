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

func provideRegistry() ports.Registry {
	r := embed_registry.New()
	register.Builtins(r)
	return r
}

func hasVerboseFlag(args []string) bool {
	for _, a := range args {
		if a == "--verbose" || a == "-v" { // simple pre-scan; Cobra will parse formally later
			return true
		}
		// Handle combined short flags like -vv or -vfoo (we only care about presence of 'v' by itself)
		if strings.HasPrefix(a, "-") && strings.Contains(a, "v") && !strings.HasPrefix(a, "--") {
			return true
		}
	}
	return false
}

func main() {
	verbose := hasVerboseFlag(os.Args[1:])
	platformlog.Init(verbose)

	fxLoggerOpt := fx.WithLogger(func() fxevent.Logger {
		if verbose {
			return &fxevent.SlogLogger{Logger: platformlog.L()}
		}
		return fxevent.NopLogger
	})

	app := fx.New(
		fxLoggerOpt,
		fx.Provide(
			provideRegistry,
		),
		fx.Invoke(func(reg ports.Registry) error {
			root := cli.NewRootCmd(reg)
			if err := root.Execute(); err != nil {
				return err
			}
			return nil
		}),
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		os.Exit(1)
	}
	// Stop immediately after the command completes (Start will return once Invoke returns)
	_ = app.Stop(ctx)
}
