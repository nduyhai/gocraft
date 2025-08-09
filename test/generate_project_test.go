package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/you/cleanctl/internal/adapters/outbound/fs/oswriter"
	"github.com/you/cleanctl/internal/adapters/outbound/rendering/texttmpl"
	embedrepo "github.com/you/cleanctl/internal/adapters/outbound/templates/embed_repo"
	"github.com/you/cleanctl/internal/core/entity"
	"github.com/you/cleanctl/internal/core/usecase"
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
	// Verify a file exists
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		t.Fatalf("expected go.mod, got err: %v", err)
	}
}
