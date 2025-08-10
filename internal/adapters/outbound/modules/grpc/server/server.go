package grpcservermodule

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/nduyhai/gocraft/internal/core/ports"
)

// Module implements ports.Module for adding a gRPC server into an existing or new project.
//
// Name:     grpc:server
// Requires: platform:base (for project structure and Fx)
// Conflicts: none
//
// This module adds:
// - internal/adapters/inbound/grpc/server/ (server wiring)
// - Registers google.golang.org/grpc/health checking service by default
// - Minimal config via Viper with default addr ":9090"

type Module struct{}

func New() Module { return Module{} }

func (Module) Name() string    { return "grpc:server" }
func (Module) Label() string   { return "gRPC Server" }
func (Module) Version() string { return "0.1.0" }
func (Module) Summary() string {
	return "Adds a gRPC server with default Google health check service"
}
func (Module) Tags() []string { return []string{"grpc", "server"} }

func (Module) Requires() []string  { return []string{"platform:base"} }
func (Module) Conflicts() []string { return nil }

func (Module) Applies(ctx ports.Ctx) bool { return true }

func (Module) Apply(ctx ports.Ctx) error {
	// Try to add required dependencies if a GoMod editor is available
	if gm := ctx.GoMod(); gm != nil {
		_ = gm.Add("google.golang.org/grpc", "v1.63.2")
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
	tpl := ports.Template{Name: "grpc:server", Files: tfiles}
	files, err := ctx.Renderer().Render(tpl, ctx.Values())
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	if err := ctx.FS().WriteAll(ctx.ProjectRoot(), files); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// Update the generated project's DI root to append the gRPC module via AdaptersModule editor.
	if ed := ctx.AdaptersModule(); ed != nil {
		vals := ctx.Values()
		modPath, _ := vals["Module"].(string)
		if modPath == "" {
			return fmt.Errorf("module path missing in context")
		}
		if err := ed.Ensure("grpcserver", modPath+"/internal/adapters/inbound/grpc/server", "grpcserver.Module()"); err != nil {
			return fmt.Errorf("update di root: %w", err)
		}
	}
	return nil
}
