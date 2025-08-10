package cli

import (
	"fmt"
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

func newNewCmd(reg ports.Registry) *cobra.Command {
	var (
		module string
		with   []string
		set    []string
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
			cfgEditor := configfileeditor.New(target)
			vals := map[string]any{"Name": name, "Module": module}
			if len(set) > 0 {
				mergeSetsInto(vals, set)
			}
			ctx := contextimpl.New(
				target,
				writer,
				renderer,
				gomod,
				adaptersEditor,
				cfgEditor,
				vals,
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
	cmd.Flags().StringSliceVar(&set, "set", nil, "Set template values (key=value). Supports dot paths, e.g., --set gorm.driver=postgres")
	// Template (-t) and output (-o) flags are no longer needed; default template is platform:base via modules and output is ./<name>
	return cmd
}

// mergeSetsInto parses key=value pairs and merges them into vals map with dot-path nesting
func mergeSetsInto(vals map[string]any, sets []string) {
	for _, kv := range sets {
		if kv == "" {
			continue
		}
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := parts[1]
		if key == "" {
			continue
		}
		setNested(vals, strings.Split(key, "."), val)
	}
}

func setNested(m map[string]any, path []string, value any) {
	cur := m
	for i, p := range path {
		if i == len(path)-1 {
			cur[p] = value
			return
		}
		next, ok := cur[p]
		if !ok {
			nm := make(map[string]any)
			cur[p] = nm
			cur = nm
			continue
		}
		nm, ok := next.(map[string]any)
		if !ok {
			nm = make(map[string]any)
			cur[p] = nm
		}
		cur = nm
	}
}
