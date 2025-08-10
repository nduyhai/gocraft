package makefile

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/nduyhai/gocraft/internal/core/ports"
)

// Module implements ports.Module for adding a Makefile to the project.
//
// Name:      feature:makefile
// Requires:  platform:base (so the base project exists and paths are consistent)
// Conflicts: none
//
// This module writes a Makefile at the project root from an embedded template.
// The template is conservative and idempotent in that it will overwrite the
// existing Makefile if present; future enhancement could add merge/skip logic.

type Module struct{}

func New() Module { return Module{} }

func (Module) Name() string    { return "feature:makefile" }
func (Module) Label() string   { return "Makefile (common dev tasks)" }
func (Module) Version() string { return "0.1.0" }
func (Module) Summary() string { return "Adds a Makefile with common targets (build, test, lint, run)" }
func (Module) Tags() []string  { return []string{"feature", "makefile", "devtools"} }

func (Module) Requires() []string  { return []string{"platform:base"} }
func (Module) Conflicts() []string { return nil }

func (Module) Applies(ctx ports.Ctx) bool { return true }

func (Module) Apply(ctx ports.Ctx) error {
	// Build template from embedded TemplatesFS
	sub, err := fs.Sub(TemplatesFS, "templates")
	if err != nil {
		return fmt.Errorf("sub fs: %w", err)
	}
	var tfiles []ports.TmplFile
	err = fs.WalkDir(sub, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(sub, path)
		if err != nil {
			return err
		}
		clean := strings.TrimPrefix(path, "./")
		tfiles = append(tfiles, ports.TmplFile{Path: clean, Content: string(b)})
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk: %w", err)
	}
	tpl := ports.Template{Name: "feature:makefile", Files: tfiles}
	files, err := ctx.Renderer().Render(tpl, ctx.Values())
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	if err := ctx.FS().WriteAll(ctx.ProjectRoot(), files); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// Defaults returns no defaults for this feature module.
func (Module) Defaults() map[string]any { return nil }
