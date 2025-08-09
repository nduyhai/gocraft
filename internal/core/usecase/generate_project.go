package usecase

import (
	"context"
	"fmt"

	"github.com/nduyhai/go-clean-arch-starter/internal/core/entity"
	"github.com/nduyhai/go-clean-arch-starter/internal/core/ports"
)

// GenerateProject orchestrates template repo -> render -> write -> hooks.
type GenerateProject struct {
	Templates ports.TemplateRepo
	Renderer  ports.Renderer
	Writer    ports.FSWriter
	Hook      ports.PostHook
}

func (uc GenerateProject) Execute(ctx context.Context, p entity.Project) error {
	if p.Template == "" {
		p.Template = "basic"
	}
	// Load template
	tpl, err := uc.Templates.Load(p.Template)
	if err != nil {
		return fmt.Errorf("load template: %w", err)
	}
	// Render
	files, err := uc.Renderer.Render(tpl, map[string]any{
		"Name":   p.Name,
		"Module": p.Module,
	})
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	// Write
	if err := uc.Writer.WriteAll(p.TargetDir, files); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	// Post hook (best-effort)
	if uc.Hook != nil {
		if err := uc.Hook.Run(p.TargetDir); err != nil {
			return fmt.Errorf("posthook: %w", err)
		}
	}
	return nil
}
