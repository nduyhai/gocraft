package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nduyhai/gocraft/internal/core/ports"
	"github.com/nduyhai/gocraft/internal/core/usecase"
	"github.com/spf13/cobra"
)

func newNewCmd(reg ports.Registry) *cobra.Command {
	var (
		module string
		with   []string
	)

	cmd := &cobra.Command{
		Use:   "new [project]",
		Short: "Generate a new project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if module == "" {
				module = fmt.Sprintf("github.com/you/%s", name)
			}
			target := filepath.Join(".", name)

			// Build module context with common collaborators
			ctx := newCtx(target, name, module)

			// Use usecase to apply module(s) with injected registry
			uc := usecase.ApplyModules{Registry: reg}
			mods := append([]string{"platform:base"}, with...)
			if err := uc.Execute(ctx, mods...); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Project generated at %s\n", target)
			if len(with) > 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Applied modules: %s\n", strings.Join(with, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&module, "module", "m", "", "Go module path (default: github.com/you/<name>)")
	cmd.Flags().StringSliceVar(&with, "with", nil, "Additional modules to apply (e.g. --with http:gin)")
	// Template (-t) and output (-o) flags are no longer needed; default template is platform:base via modules and output is ./<name>
	return cmd
}
