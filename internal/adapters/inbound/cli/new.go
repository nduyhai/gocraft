package cli

import (
	"fmt"
	"path/filepath"

	contextimpl "github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/context/contextimpl"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/fs/oswriter"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/modules/register"
	regsimple "github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/registry/simple"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/rendering/texttmpl"
	"github.com/spf13/cobra"
)

func newNewCmd() *cobra.Command {
	var (
		module string
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

			// Build registry and register built-ins, then apply platform:base
			r := regsimple.New()
			register.Builtins(r)
			if err := r.Apply(ctx, "platform:base"); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Project generated at %s\n", target)
			return nil
		},
	}

	cmd.Flags().StringVarP(&module, "module", "m", "", "Go module path (default: github.com/you/<name>)")
	// Template (-t) and output (-o) flags are no longer needed; default template is platform:base via modules and output is ./<name>
	return cmd
}
