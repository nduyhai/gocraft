package ginmodule

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/nduyhai/gocraft/internal/core/ports"
	"gopkg.in/yaml.v3"
)

// Module implements ports.Module for adding a Gin HTTP server into an existing or new project.
//
// Name:     http:gin
// Requires: platform:base (for project structure and Fx)
// Conflicts: none
//
// This module adds:
// - internal/adapters/http/gin/ (server wiring, middlewares)
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
func (Module) Conflicts() []string { return []string{"http:chi"} }

func (Module) Applies(ctx ports.Ctx) bool { return true }

func (Module) Apply(ctx ports.Ctx) error {
	// Try to add required dependencies if a GoMod editor is available
	if gm := ctx.GoMod(); gm != nil {
		_ = gm.Add("github.com/gin-gonic/gin", "v1.10.0")
		_ = gm.Add("github.com/google/uuid", "v1.6.0")
		_ = gm.Add("github.com/prometheus/client_golang", "v1.19.1")
		_ = gm.Add("github.com/spf13/viper", "v1.20.1")
		_ = gm.Add("go.uber.org/fx", "v1.24.0")
	}
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
	tpl := ports.Template{Name: "http:gin", Files: tfiles}
	files, err := ctx.Renderer().Render(tpl, ctx.Values())
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	if err := ctx.FS().WriteAll(ctx.ProjectRoot(), files); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// Update the generated project's DI root to append the Gin module via AdaptersModule editor.
	if ed := ctx.AdaptersModule(); ed != nil {
		vals := ctx.Values()
		modPath, _ := vals["Module"].(string)
		if modPath == "" {
			return fmt.Errorf("module path missing in context")
		}
		if err := ed.Ensure("httpgin", modPath+"/internal/adapters/inbound/http/gin", "httpgin.Module()"); err != nil {
			return fmt.Errorf("update di root: %w", err)
		}
	}
	return nil
}

// Defaults implements ports.Module.Defaults to provide default configuration.
func (Module) Defaults() map[string]any {
	// Attempt to load from embedded defaults template
	if b, err := fs.ReadFile(TemplatesFS, "templates/config/defaults.yml.tmpl"); err == nil {
		var m map[string]any
		if yaml.Unmarshal(b, &m) == nil && m != nil {
			return m
		}
	}
	// Fallback to hardcoded defaults
	return map[string]any{
		"server": map[string]any{
			"http": map[string]any{
				"addr": ":8080",
			},
		},
	}
}
