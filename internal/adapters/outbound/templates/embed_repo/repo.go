package embed_repo

import (
	"fmt"
	"io/fs"
	"strings"

	ginmodule "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/http/gin"
	base "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/platform/base"
	"github.com/nduyhai/gocraft/internal/core/ports"
)

// New returns an embedded template repository with available templates.
func New() ports.TemplateRepo { return &repo{names: []string{"basic", "http:gin"}} }

type repo struct{ names []string }

func (r *repo) Names() []string { return append([]string(nil), r.names...) }

func (r *repo) Load(name string) (ports.Template, error) {
	var (
		fsys fs.FS
	)
	switch name {
	case "basic":
		fsys = base.TemplatesFS
	case "http:gin":
		fsys = ginmodule.TemplatesFS
	default:
		return ports.Template{}, fmt.Errorf("unknown template: %s", name)
	}
	// We want paths relative to the root of the templates dir
	sub, err := fs.Sub(fsys, "templates")
	if err != nil {
		return ports.Template{}, err
	}
	var files []ports.TmplFile
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
		files = append(files, ports.TmplFile{Path: clean, Content: string(b)})
		return nil
	})
	if err != nil {
		return ports.Template{}, err
	}
	return ports.Template{Name: name, Files: files}, nil
}
