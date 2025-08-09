package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/fs/oswriter"
	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/rendering/texttmpl"
	embedrepo "github.com/nduyhai/go-clean-arch-starter/internal/adapters/outbound/templates/embed_repo"
	"github.com/nduyhai/go-clean-arch-starter/internal/core/entity"
	"github.com/nduyhai/go-clean-arch-starter/internal/core/usecase"
)

func TestGenerateProject_Basic(t *testing.T) {
	dir := t.TempDir()
	repo := embedrepo.New()
	renderer := texttmpl.New()
	writer := oswriter.New()
	uc := usecase.GenerateProject{Templates: repo, Renderer: renderer, Writer: writer}

	p := entity.Project{
		Name:      "demo",
		Module:    "github.com/you/demo",
		Template:  "basic",
		TargetDir: dir,
	}

	if err := uc.Execute(context.Background(), p); err != nil {
		t.Fatalf("execute: %v", err)
	}
	// Verify go.mod exists
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		t.Fatalf("expected go.mod, got err: %v", err)
	}
	// Verify cmd/<name>/main.go exists (ensure __name__ dir is embedded)
	if _, err := os.Stat(filepath.Join(dir, "cmd", "demo", "main.go")); err != nil {
		t.Fatalf("expected cmd/demo/main.go, got err: %v", err)
	}
	// Verify .gitkeep-based directories are created by ensuring .gitkeep is written
	if _, err := os.Stat(filepath.Join(dir, "internal", "core", "entity", ".gitkeep")); err != nil {
		t.Fatalf("expected internal/core/entity/.gitkeep, got err: %v", err)
	}
}
