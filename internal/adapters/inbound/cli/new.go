package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nduyhai/gocraft/internal/adapters/outbound/context/contextimpl"
	amfileeditor "github.com/nduyhai/gocraft/internal/adapters/outbound/editors/adaptersmodule/fileeditor"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/fs/oswriter"
	gomodfileeditor "github.com/nduyhai/gocraft/internal/adapters/outbound/gomod/fileeditor"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/rendering/texttmpl"
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

			// Prepare outbound collaborators
			renderer := texttmpl.New()
			writer := oswriter.New()

			// Build module context
			// Build GoMod editor bound to target directory
			gomod := gomodfileeditor.New(target)
			adaptersEditor := amfileeditor.New(target)
			ctx := contextimpl.New(
				target,
				writer,
				renderer,
				gomod,
				adaptersEditor,
				map[string]any{"Name": name, "Module": module},
			)

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
