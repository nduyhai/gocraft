package gitignore

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/nduyhai/gocraft/internal/core/ports"
)

// Module implements ports.Module for adding a .gitignore file to the project root.
//
// Name:      feature:gitignore
// Requires:  platform:base (ensure base project layout exists)
// Conflicts: none
//
// This module writes a .gitignore at the project root from an embedded template.
// It is idempotent with respect to file generation (will overwrite if exists).

type Module struct{}

func New() Module { return Module{} }

func (Module) Name() string    { return "feature:gitignore" }
func (Module) Label() string   { return ".gitignore (Go project defaults)" }
func (Module) Version() string { return "0.1.0" }
func (Module) Summary() string { return "Adds a .gitignore suited for Go projects" }
func (Module) Tags() []string  { return []string{"feature", "git", "ignore"} }

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
	tpl := ports.Template{Name: "feature:gitignore", Files: tfiles}
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
