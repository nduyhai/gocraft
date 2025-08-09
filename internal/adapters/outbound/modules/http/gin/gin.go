package ginmodule

import (
	"fmt"

	embedrepo "github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/templates/embed_repo"
	"github.com/nduyhai/go-clean-arch-starter/internal/core/ports"
)

// Module implements ports.Module for adding a Gin HTTP server into an existing or new project.
//
// Name:     http:gin
// Requires: platform:base (for project structure and Fx)
// Conflicts: none
//
// This module adds:
// - internal/adapters/http/gin/ (server wiring, middlewares)
// - Self-registration into adapters.Register(...) so no file overwrite is needed
// - Default routes: /healthz and /metrics
// - Middlewares: recovery, request ID, simple logging, permissive CORS
// - Minimal config via Viper with default addr ":8080" and PORT env override

type Module struct{}

func New() Module { return Module{} }

func (Module) Name() string    { return "http:gin" }
func (Module) Label() string   { return "HTTP Gin Server" }
func (Module) Version() string { return "0.1.0" }
func (Module) Summary() string {
	return "Adds a Gin HTTP server with /metrics, /healthz and common middlewares"
}
func (Module) Tags() []string { return []string{"http", "gin", "server"} }

func (Module) Requires() []string  { return []string{"platform:base"} }
func (Module) Conflicts() []string { return nil }

func (Module) Applies(ctx ports.Ctx) bool { return true }

func (Module) Apply(ctx ports.Ctx) error {
	repo := embedrepo.New()
	tpl, err := repo.Load("http:gin")
	if err != nil {
		return fmt.Errorf("load template: %w", err)
	}
	files, err := ctx.Renderer().Render(tpl, ctx.Values())
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	if err := ctx.FS().WriteAll(ctx.ProjectRoot(), files); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
