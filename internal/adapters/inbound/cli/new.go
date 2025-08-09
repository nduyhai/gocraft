package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/you/cleanctl/internal/adapters/outbound/fs/oswriter"
	"github.com/you/cleanctl/internal/adapters/outbound/hooks/exec"
	"github.com/you/cleanctl/internal/adapters/outbound/rendering/texttmpl"
	"github.com/you/cleanctl/internal/adapters/outbound/templates/embed_repo"
	"github.com/you/cleanctl/internal/core/entity"
	"github.com/you/cleanctl/internal/core/usecase"
)

func newNewCmd() *cobra.Command {
	var (
		module   string
		template string
		target   string
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
			if target == "" {
				target = filepath.Join(".", name)
			}

			repo := embed_repo.New()
			renderer := texttmpl.New()
			writer := oswriter.New()
			hook := exec.New()
			uc := usecase.GenerateProject{Templates: repo, Renderer: renderer, Writer: writer, Hook: hook}

			p := entity.Project{
				Name:      name,
				Module:    module,
				Template:  template,
				Options:   map[string]string{},
				TargetDir: target,
			}

			if err := uc.Execute(context.Background(), p); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Project generated at %s\n", target)
			return nil
		},
	}

	cmd.Flags().StringVarP(&module, "module", "m", "", "Go module path (default: github.com/you/<name>)")
	cmd.Flags().StringVarP(&template, "template", "t", "basic", "Template name")
	cmd.Flags().StringVarP(&target, "output", "o", "", "Target directory (default: ./<name>)")
	return cmd
}
