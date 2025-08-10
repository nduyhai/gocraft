package cli

import (
	"os"

	"github.com/nduyhai/gocraft/internal/core/ports"
	"github.com/nduyhai/gocraft/pkg/version"
	"github.com/spf13/cobra"
)

func NewRootCmd(reg ports.Registry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gocraft",
		Short: "Project generator following clean architecture",
		Long:  "gocraft generates a Go project from embedded templates using a clean architecture.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.Version = version.Version

	cmd.AddCommand(newNewCmd(reg))
	cmd.AddCommand(newListCmd(reg))
	cmd.AddCommand(newAddCmd(reg))
	cmd.AddCommand(newCompletionCmd())

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	cmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")

	cmd.SetErrPrefix("error: ")
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	return cmd
}
