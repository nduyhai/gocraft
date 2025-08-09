package texttmpl

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nduyhai/gocraft/internal/core/entity"
	"github.com/nduyhai/gocraft/internal/core/ports"
)

type Renderer struct{}

func New() *Renderer { return &Renderer{} }

func (Renderer) Render(tpl ports.Template, ctx any) ([]entity.File, error) {
	var out []entity.File
	for _, f := range tpl.Files {
		// Render file content
		content, err := renderString(f.Content, ctx)
		if err != nil {
			return nil, err
		}
		// Path token rewrites and .tmpl stripping
		path := applyPathTokens(filepath.FromSlash(f.Path), ctx)
		if strings.HasSuffix(path, ".tmpl") {
			path = strings.TrimSuffix(path, ".tmpl")
		}
		out = append(out, entity.File{Path: path, Content: []byte(content), Mode: 0o644})
	}
	return out, nil
}

func renderString(tmpl string, data any) (string, error) {
	funcs := template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"kebab": toKebab,
		"snake": toSnake,
	}
	t, err := template.New("").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return strings.ReplaceAll(buf.String(), "\r\n", "\n"), nil
}

// applyPathTokens replaces path tokens like __name__ and __module__ using context values.
func applyPathTokens(path string, ctx any) string {
	name := getCtxString(ctx, "Name")
	module := getCtxString(ctx, "Module")
	if name != "" {
		path = strings.ReplaceAll(path, "__name__", name)
	}
	if module != "" {
		path = strings.ReplaceAll(path, "__module__", module)
	}
	// Also support kebab/snake variants if template authors use them in paths
	if name != "" {
		path = strings.ReplaceAll(path, "__name-kebab__", toKebab(name))
		path = strings.ReplaceAll(path, "__name-snake__", toSnake(name))
	}
	return path
}

func getCtxString(ctx any, key string) string {
	m, ok := ctx.(map[string]any)
	if !ok {
		return ""
	}
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func toKebab(s string) string {
	parts := splitWords(s)
	for i := range parts {
		parts[i] = strings.ToLower(parts[i])
	}
	return strings.Join(parts, "-")
}

func toSnake(s string) string {
	parts := splitWords(s)
	for i := range parts {
		parts[i] = strings.ToLower(parts[i])
	}
	return strings.Join(parts, "_")
}

func splitWords(s string) []string {
	if s == "" {
		return nil
	}
	runes := []rune(s)
	var parts []string
	start := 0
	for i := 1; i < len(runes); i++ {
		prev := runes[i-1]
		cur := runes[i]
		// boundary if letter case changes (aB or ABc), or non-alnum
		if isBoundary(prev, cur) {
			parts = append(parts, string(runes[start:i]))
			start = i
		}
	}
	parts = append(parts, string(runes[start:]))
	return parts
}

func isBoundary(prev, cur rune) bool {
	if !isAlphaNum(prev) || !isAlphaNum(cur) {
		return true
	}
	// lower -> upper
	if isLower(prev) && isUpper(cur) {
		return true
	}
	// letter -> digit or digit -> letter
	if isDigit(prev) != isDigit(cur) {
		return true
	}
	return false
}

func isAlphaNum(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}
func isLower(r rune) bool { return r >= 'a' && r <= 'z' }
func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }
func isDigit(r rune) bool { return r >= '0' && r <= '9' }
