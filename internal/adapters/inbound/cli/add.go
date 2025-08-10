package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nduyhai/gocraft/internal/core/ports"
	"github.com/nduyhai/gocraft/internal/core/usecase"
	"github.com/spf13/cobra"
)

// newAddCmd creates the `add` command which applies one or more modules to the current project directory.
func newAddCmd(reg ports.Registry) *cobra.Command {
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

			// Build context with common collaborators
			ctx := newCtx(cwd, name, modulePath)

			// Use usecase to apply modules with injected registry
			uc := usecase.ApplyModules{Registry: reg}
			if err := uc.Execute(ctx, args...); err != nil {
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
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("module path not found in go.mod")
}
