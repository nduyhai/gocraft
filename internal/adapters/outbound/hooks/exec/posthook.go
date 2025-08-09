package exec

import (
	"os/exec"
)

type Hook struct{}

func New() *Hook { return &Hook{} }

// Run executes a few best-effort commands: git init and go mod tidy.
func (Hook) Run(targetDir string) error {
	// ignore errors intentionally (best effort)
	_ = run(targetDir, "git", "init")
	_ = run(targetDir, "git", "add", ".")
	_ = run(targetDir, "git", "commit", "-m", "chore: initial commit")
	_ = run(targetDir, "go", "mod", "tidy")
	return nil
}

func run(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}
