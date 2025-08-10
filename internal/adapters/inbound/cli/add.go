package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	configfileeditor "github.com/nduyhai/gocraft/internal/adapters/outbound/config/fileeditor"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/context/contextimpl"
	amfileeditor "github.com/nduyhai/gocraft/internal/adapters/outbound/di/fileeditor"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/fs/oswriter"
	gomodfileeditor "github.com/nduyhai/gocraft/internal/adapters/outbound/gomod/fileeditor"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/rendering/texttmpl"
	"github.com/nduyhai/gocraft/internal/core/ports"
	"github.com/nduyhai/gocraft/internal/core/usecase"
	"github.com/spf13/cobra"
)

// newAddCmd creates the `add` command which applies one or more modules to the current project directory.
func newAddCmd(reg ports.Registry) *cobra.Command {
	var set []string
	cmd := &cobra.Command{
		Use:   "add <module>...",
		Short: "Apply module(s) to the current project",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("getwd: %w", err)
			}

			// Determine Name and Module for rendering context
			name := filepath.Base(cwd)
			modulePath, err := readModulePath(filepath.Join(cwd, "go.mod"))
			if err != nil {
				// Fallback to a sensible default if go.mod is missing
				modulePath = fmt.Sprintf("github.com/you/%s", name)
			}

			// Outbound collaborators
			renderer := texttmpl.New()
			writer := oswriter.New()
			gomod := gomodfileeditor.New(cwd)
			adaptersEditor := amfileeditor.New(cwd)
			cfgEditor := configfileeditor.New(cwd)

			// Build context
			vals := map[string]any{
				"Name":   name,
				"Module": modulePath,
			}
			if len(set) > 0 {
				mergeSetsInto(vals, set)
			}
			ctx := contextimpl.New(cwd, writer, renderer, gomod, adaptersEditor, cfgEditor, vals)

			// Use usecase to apply modules with injected registry
			uc := usecase.ApplyModules{Registry: reg}
			if err := uc.Execute(ctx, args...); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Applied modules: %s\n", strings.Join(args, ", "))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&set, "set", nil, "Set template values (key=value). Supports dot paths, e.g., --set gorm.driver=postgres")
	return cmd
}

// readModulePath reads the module path from a go.mod file. Returns an error if the file
// cannot be read or the module line is not found.
func readModulePath(goModPath string) (string, error) {
	f, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	Scanner := bufio.NewScanner(f)
	for Scanner.Scan() {
		line := strings.TrimSpace(Scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	if err := Scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("module path not found in go.mod")
}
