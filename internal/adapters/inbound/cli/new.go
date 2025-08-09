package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/context/contextimpl"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/fs/oswriter"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/modules/register"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/registry/embed_registry"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/rendering/texttmpl"
	"github.com/nduyhai/go-clean-arch-starter/internal/core/usecase"
	"github.com/spf13/cobra"
)

func newNewCmd() *cobra.Command {
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

			// Prepare outbound collaborators
			renderer := texttmpl.New()
			writer := oswriter.New()

			// Build module context
			ctx := contextimpl.New(
				target,
				writer,
				renderer,
				nil, // GoModEditor not needed for base generation yet
				map[string]any{"Name": name, "Module": module},
			)

			// Build registry and register built-ins
			r := embed_registry.New()
			register.Builtins(r)

			// Use usecase to apply module(s)
			uc := usecase.ApplyModules{Registry: r}
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
