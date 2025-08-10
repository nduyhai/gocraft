package cli

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/nduyhai/gocraft/internal/core/ports"
	"github.com/nduyhai/gocraft/internal/core/usecase"
	"github.com/spf13/cobra"
)

func newListCmd(reg ports.Registry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use the injected registry to list modules
			uc := usecase.ListModules{Registry: reg}
			mods := uc.Execute()
			if len(mods) == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No modules available")
				return nil
			}

			// Sort by Name for stable output (registry already preserves order, but sort for clarity)
			sort.Slice(mods, func(i, j int) bool { return mods[i].Name() < mods[j].Name() })

			// Print table using tabwriter for aligned columns
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "NAME\tVERSION\tLABEL\tTAGS\tSUMMARY")
			for _, m := range mods {
				name := m.Name()
				version := m.Version()
				label := m.Label()
				tags := strings.Join(m.Tags(), ",")
				summary := m.Summary()
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", name, version, label, tags, summary)
			}
			if err := w.Flush(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
