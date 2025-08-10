package fileeditor_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	amedit "github.com/nduyhai/gocraft/internal/adapters/outbound/di/fileeditor"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestEditor_Ensure_AST_Idempotent(t *testing.T) {
	dir := t.TempDir()
	rootGo := filepath.Join(dir, "internal", "platform", "di", "root.go")
	write(t, rootGo, `package di

import (
	"go.uber.org/fx"
	"example.com/app/internal/platform/env"
	"example.com/app/internal/platform/logger"
)

func Root() fx.Option {
	return fx.Options(
		env.Module(),
		logger.Module(),
	)
}
`)

	ed := amedit.New(dir)
	alias := "httpgin"
	imp := "example.com/app/internal/adapters/inbound/http/gin"
	optexpr := "httpgin.Module()"

	if err := ed.Ensure(alias, imp, optexpr); err != nil {
		t.Fatalf("Ensure 1: %v", err)
	}
	b1, _ := os.ReadFile(rootGo)
	got1 := string(b1)
	if !strings.Contains(got1, "\t"+alias+" \""+imp+"\"") {
		t.Fatalf("import not found after first ensure. got=\n%s", got1)
	}
	if !strings.Contains(got1, optexpr) {
		t.Fatalf("option expr not found after first ensure. got=\n%s", got1)
	}

	// Run again should be idempotent (no duplicate occurrences)
	if err := ed.Ensure(alias, imp, optexpr); err != nil {
		t.Fatalf("Ensure 2: %v", err)
	}
	b2, _ := os.ReadFile(rootGo)
	got2 := string(b2)
	if c := strings.Count(got2, "\t"+alias+" \""+imp+"\""); c != 1 {
		t.Fatalf("import occurrence count = %d, want 1. file=\n%s", c, got2)
	}
	if c := strings.Count(got2, optexpr); c != 1 {
		t.Fatalf("option expr occurrence count = %d, want 1. file=\n%s", c, got2)
	}
}
