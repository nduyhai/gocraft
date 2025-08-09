package embed_repo

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nduyhai/go-clean-arch-starter/internal/core/ports"
)

//go:embed templates/** templates/basic/cmd/__name__/**
//go:embed templates/basic/internal/adapters/inbound/.gitkeep
//go:embed templates/basic/internal/adapters/outbound/.gitkeep
//go:embed templates/basic/internal/core/entity/.gitkeep
//go:embed templates/basic/internal/core/ports/.gitkeep
//go:embed templates/basic/internal/core/usecase/.gitkeep
var filesFS embed.FS

// New returns an embedded template repository with available templates.
func New() ports.TemplateRepo { return &repo{names: []string{"basic", "http:gin"}} }

type repo struct{ names []string }

func (r *repo) Names() []string { return append([]string(nil), r.names...) }

func (r *repo) Load(name string) (ports.Template, error) {
	// Map template name to its base directory under templates/
	var base string
	switch name {
	case "basic":
		base = "templates/basic"
	case "http:gin":
		base = "templates/http/gin"
	default:
		return ports.Template{}, fmt.Errorf("unknown template: %s", name)
	}
	var files []ports.TmplFile
	err := fs.WalkDir(filesFS, base, func(path string, d fs.DirEntry, err error) error {
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
		rel := strings.TrimPrefix(path, base+"/")
		files = append(files, ports.TmplFile{Path: rel, Content: string(b)})
		return nil
	})
	if err != nil {
		return ports.Template{}, err
	}
	return ports.Template{Name: name, Files: files}, nil
}
