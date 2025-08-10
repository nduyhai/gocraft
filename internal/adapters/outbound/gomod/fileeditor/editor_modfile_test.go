package fileeditor_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	gmedit "github.com/nduyhai/gocraft/internal/adapters/outbound/gomod/fileeditor"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestGoModEditor_Add_Replace_Idempotent(t *testing.T) {
	dir := t.TempDir()
	gomod := filepath.Join(dir, "go.mod")
	writeFile(t, gomod, `module example.com/app

go 1.22

require (
	go.uber.org/fx v1.22.0
)
`)

	ed := gmedit.New(dir)
	// Add new require
	if err := ed.Add("github.com/spf13/cobra", "v1.8.1"); err != nil {
		t.Fatalf("Add 1: %v", err)
	}
	b1, _ := os.ReadFile(gomod)
	got1 := string(b1)
	if !strings.Contains(got1, "github.com/spf13/cobra v1.8.1") {
		t.Fatalf("missing require after add: %s", got1)
	}
	// Idempotent
	if err := ed.Add("github.com/spf13/cobra", "v1.8.1"); err != nil {
		t.Fatalf("Add 2: %v", err)
	}
	b2, _ := os.ReadFile(gomod)
	if c := strings.Count(string(b2), "github.com/spf13/cobra v1.8.1"); c != 1 {
		t.Fatalf("duplicate require count=%d", c)
	}

	// Replace update
	if err := ed.Replace("go.uber.org/fx", "../fxfork"); err != nil {
		t.Fatalf("Replace: %v", err)
	}
	b3, _ := os.ReadFile(gomod)
	got3 := string(b3)
	if !strings.Contains(got3, "replace go.uber.org/fx => ../fxfork") {
		t.Fatalf("missing replace: %s", got3)
	}
	// Idempotent update
	if err := ed.Replace("go.uber.org/fx", "../fxfork"); err != nil {
		t.Fatalf("Replace 2: %v", err)
	}
	b4, _ := os.ReadFile(gomod)
	if c := strings.Count(string(b4), "replace go.uber.org/fx => ../fxfork"); c != 1 {
		t.Fatalf("duplicate replace count=%d", c)
	}
}
