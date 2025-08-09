package base

import (
	"fmt"

	embedrepo "github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/templates/embed_repo"
	"github.com/nduyhai/go-clean-arch-starter/internal/core/ports"
)

// Module implements ports.Module for the base platform.
// It generates a base project layout (Fx + Viper, logger, DI root, basic core/adapters skeleton)
// using the embedded "basic" template.
//
// Name:     platform:base
// Requires: none
// Conflicts:none

type Module struct{}

func New() Module { return Module{} }

func (Module) Name() string    { return "platform:base" }
func (Module) Label() string   { return "Platform Base (Fx + Viper)" }
func (Module) Version() string { return "0.1.0" }
func (Module) Summary() string {
	return "Generates base clean-arch project (Fx DI, Viper config, logger, DI root)"
}
func (Module) Tags() []string { return []string{"platform", "base", "fx", "viper"} }

func (Module) Requires() []string  { return nil }
func (Module) Conflicts() []string { return nil }

// Applies returns true if we should apply the module in the given context.
// For now, always true; in future we could detect if files already exist to avoid overwrite.
func (Module) Applies(ctx ports.Ctx) bool { return true }

func (Module) Apply(ctx ports.Ctx) error {
	// Load embedded template
	repo := embedrepo.New()
	tpl, err := repo.Load("basic")
	if err != nil {
		return fmt.Errorf("load template: %w", err)
	}
	// Render with provided values (.Name, .Module)
	files, err := ctx.Renderer().Render(tpl, ctx.Values())
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	// Write to project root
	if err := ctx.FS().WriteAll(ctx.ProjectRoot(), files); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
