package cli

import (
	"fmt"

	"github.com/nduyhai/gocraft/pkg/version"
	"github.com/spf13/cobra"
)

// newVersionCmd creates the `version` command which prints the gocraft version.
func newVersionCmd() *cobra.Command {
	var short bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show gocraft version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if short {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), version.Version)
				return nil
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "gocraft version %s\n", version.Version)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&short, "short", "s", false, "Print just the version number")
	return cmd
}
