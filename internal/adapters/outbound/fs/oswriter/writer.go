package oswriter

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nduyhai/gocraft/internal/core/entity"
)

type Writer struct{}

func New() *Writer { return &Writer{} }

func (Writer) WriteAll(root string, files []entity.File) error {
	for _, f := range files {
		path := filepath.Join(root, f.Path)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
		if err := writeFile(path, f.Content, f.Mode); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(path string, content []byte, mode fs.FileMode) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file exists: %s", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	return os.WriteFile(path, content, mode)
}
