package embed_repo

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/you/cleanctl/internal/core/ports"
)

//go:embed templates/**
var filesFS embed.FS

// New returns an embedded template repository with a single template named "basic".
func New() ports.TemplateRepo { return &repo{names: []string{"basic"}} }

type repo struct{ names []string }

func (r *repo) Names() []string { return append([]string(nil), r.names...) }

func (r *repo) Load(name string) (ports.Template, error) {
	if name != "basic" {
		return ports.Template{}, fmt.Errorf("unknown template: %s", name)
	}
	var files []ports.TmplFile
	err := fs.WalkDir(filesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(filesFS, path)
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(path, "templates/basic/")
		files = append(files, ports.TmplFile{Path: rel, Content: string(b)})
		return nil
	})
	if err != nil {
		return ports.Template{}, err
	}
	return ports.Template{Name: name, Files: files}, nil
}
