package log

import (
	"log/slog"
	"os"
)

var defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

// Init configures the global logger. If verbose is true, sets level to Debug, otherwise Info.
func Init(verbose bool) {
	lvl := slog.LevelInfo
	if verbose {
		lvl = slog.LevelDebug
	}
	defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))
}

func L() *slog.Logger { return defaultLogger }
