package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	contextimpl "github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/context/contextimpl"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/fs/oswriter"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/modules/register"
	regsimple "github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/registry/simple"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/rendering/texttmpl"
	"github.com/spf13/cobra"
)

// newAddCmd creates the `add` command which applies one or more modules to the current project directory.
func newAddCmd() *cobra.Command {
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

			// Build context
			ctx := contextimpl.New(cwd, writer, renderer, nil, map[string]any{
				"Name":   name,
				"Module": modulePath,
			})

			// Registry and built-ins
			r := regsimple.New()
			register.Builtins(r)

			if err := r.Apply(ctx, args...); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Applied modules: %s\n", strings.Join(args, ", "))
			return nil
		},
	}
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
