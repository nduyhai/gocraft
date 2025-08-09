package cli

import (
	"os"

	"github.com/nduyhai/go-clean-arch-starter/pkg/version"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanctl",
		Short: "Project generator following clean architecture",
		Long:  "cleanctl generates a Go project from embedded templates using a clean architecture.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.Version = version.Version

	cmd.AddCommand(newNewCmd())
	cmd.AddCommand(newCompletionCmd())

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	cmd.SetErrPrefix("error: ")
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	return cmd
}
