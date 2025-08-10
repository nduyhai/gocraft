package fileeditor

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// Editor provides a go.mod editor backed by x/mod/modfile for safe edits.
// It supports adding require entries and replace directives in an idempotent way.
// Tidy remains a no-op (callers run go tooling separately).

type Editor struct {
	root string
}

func New(projectRoot string) *Editor { return &Editor{root: projectRoot} }

func (e *Editor) goModPath() string { return filepath.Join(e.root, "go.mod") }

// Add ensures the require entry exists with the specified version. If it already exists, it is left untouched.
func (e *Editor) Add(module, version string) error {
	if module == "" {
		return fmt.Errorf("module path is empty")
	}
	path := e.goModPath()
	data, err := os.ReadFile(path)
	if err != nil {
		// Missing go.mod: no-op for resilience
		return nil
	}
	mf, err := modfile.Parse(path, data, nil)
	if err != nil {
		return err
	}
	for _, r := range mf.Require {
		if r.Mod.Path == module {
			// already has a require: keep as-is to avoid breaking constraints
			return nil
		}
	}
	if err := mf.AddRequire(module, version); err != nil {
		return err
	}
	mf.Cleanup()
	formatted := modfile.Format(mf.Syntax)
	return os.WriteFile(path, formatted, 0o644)
}

// Replace adds or updates a replace directive. Versions are left empty for path-based replaces.
func (e *Editor) Replace(oldPath, newPath string) error {
	if oldPath == "" || newPath == "" {
		return fmt.Errorf("replace paths must be non-empty")
	}
	path := e.goModPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	mf, err := modfile.Parse(path, data, nil)
	if err != nil {
		return err
	}
	updated := false
	for _, rep := range mf.Replace {
		if rep.Old.Path == oldPath {
			rep.New.Path = newPath
			rep.New.Version = ""
			updated = true
		}
	}
	if !updated {
		if err := mf.AddReplace(oldPath, "", newPath, ""); err != nil {
			return err
		}
	}
	mf.Cleanup()
	formatted := modfile.Format(mf.Syntax)
	return os.WriteFile(path, formatted, 0o644)
}

// Tidy is a no-op placeholder.
func (e *Editor) Tidy() error { return nil }
