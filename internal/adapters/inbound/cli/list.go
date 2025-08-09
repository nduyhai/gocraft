package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nduyhai/gocraft/internal/adapters/outbound/modules/register"
	"github.com/nduyhai/gocraft/internal/adapters/outbound/registry/embed_registry"
	"github.com/nduyhai/gocraft/internal/core/usecase"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Build a registry and register built-in modules
			r := embed_registry.New()
			register.Builtins(r)

			// Use usecase to list modules
			uc := usecase.ListModules{Registry: r}
			mods := uc.Execute()
			if len(mods) == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No modules available")
				return nil
			}

			// Sort by Name for stable output (registry already preserves order, but sort for clarity)
			sort.Slice(mods, func(i, j int) bool { return mods[i].Name() < mods[j].Name() })

			// Print a embed_registry table-like output
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "NAME\tVERSION\tLABEL\tTAGS\tSUMMARY")
			for _, m := range mods {
				name := m.Name()
				version := m.Version()
				label := m.Label()
				tags := strings.Join(m.Tags(), ",")
				summary := m.Summary()
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\t%s\t%s\n", name, version, label, tags, summary)
			}
			return nil
		},
	}
	return cmd
}
