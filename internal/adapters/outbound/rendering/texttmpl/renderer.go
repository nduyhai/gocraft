package texttmpl

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/you/cleanctl/internal/core/entity"
	"github.com/you/cleanctl/internal/core/ports"
)

type Renderer struct{}

func New() *Renderer { return &Renderer{} }

func (Renderer) Render(tpl ports.Template, ctx any) ([]entity.File, error) {
	var out []entity.File
	for _, f := range tpl.Files {
		// directories are implicitly created by writer; we render only file contents
		content, err := renderString(f.Content, ctx)
		if err != nil {
			return nil, err
		}
		path := filepath.FromSlash(f.Path)
		if strings.HasSuffix(path, ".tmpl") {
			path = strings.TrimSuffix(path, ".tmpl")
		}
		out = append(out, entity.File{Path: path, Content: []byte(content), Mode: 0o644})
	}
	return out, nil
}

func renderString(tmpl string, data any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return strings.ReplaceAll(buf.String(), "\r\n", "\n"), nil
}
