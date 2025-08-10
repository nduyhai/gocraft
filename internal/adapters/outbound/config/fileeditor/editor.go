package configeditor

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	grpcservermodule "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/grpc/server"
	chimodule "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/http/chi"
	ginmodule "github.com/nduyhai/gocraft/internal/adapters/outbound/modules/http/gin"
	"github.com/nduyhai/gocraft/internal/core/ports"
	"gopkg.in/yaml.v3"
)

// Editor edits <root>/config/config.yml by merging module-specific default settings.
// It is idempotent: existing keys are preserved; only missing keys are added.
// It tolerates missing config file by creating it when needed.

type Editor struct{ root string }

func New(projectRoot string) *Editor { return &Editor{root: projectRoot} }

func (e *Editor) path() string { return filepath.Join(e.root, "config", "config.yml") }

// EnsureDefaultsFor merges default config for a known module into config/config.yml.
func (e *Editor) EnsureDefaultsFor(module string) error {
	defaults := defaultsFor(module)
	if defaults == nil {
		return nil
	}
	p := e.path()
	// Read existing YAML (if any)
	var current map[string]any
	b, err := os.ReadFile(p)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		current = make(map[string]any)
	} else {
		if len(b) == 0 {
			current = make(map[string]any)
		} else {
			if err := yaml.Unmarshal(b, &current); err != nil {
				// If parsing fails, do not overwrite; leave as-is (be conservative)
				return nil
			}
		}
	}
	merged := mergeMaps(current, defaults)
	out, err := yaml.Marshal(merged)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, out, 0o644)
}

var _ ports.ConfigEditor = (*Editor)(nil)

// loadDefaultsFromFS tries to read a YAML defaults template from a module's embedded FS.
// Expected path: templates/config/defaults.yml.tmpl
// Returns nil if file doesn't exist or cannot be parsed.
func loadDefaultsFromFS(fsys fs.FS) map[string]any {
	const defaultsPath = "templates/config/defaults.yml.tmpl"
	b, err := fs.ReadFile(fsys, defaultsPath)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}

// defaultsFor returns a nested map representing default config for a module.
// It first tries to load defaults from the module's embedded template FS.
// If not found, it falls back to built-in defaults to preserve existing behavior.
func defaultsFor(module string) map[string]any {
	switch module {
	case "http:gin":
		if m := loadDefaultsFromFS(ginmodule.TemplatesFS); m != nil {
			return m
		}
		return map[string]any{
			"server": map[string]any{
				"http": map[string]any{
					"addr": ":8080",
				},
			},
		}
	case "http:chi":
		if m := loadDefaultsFromFS(chimodule.TemplatesFS); m != nil {
			return m
		}
		return map[string]any{
			"server": map[string]any{
				"http": map[string]any{
					"addr": ":8080",
				},
			},
		}
	case "grpc:server":
		if m := loadDefaultsFromFS(grpcservermodule.TemplatesFS); m != nil {
			return m
		}
		return map[string]any{
			"server": map[string]any{
				"grpc": map[string]any{
					"addr":       ":9090",
					"reflection": true,
				},
			},
		}
	default:
		return nil
	}
}

// mergeMaps performs a deep merge: it sets defaults only for missing keys.
func mergeMaps(dst, def map[string]any) map[string]any {
	if dst == nil {
		dst = make(map[string]any)
	}
	for k, v := range def {
		if existing, ok := dst[k]; ok {
			// if both maps, recurse
			m1, ok1 := existing.(map[string]any)
			m2, ok2 := v.(map[string]any)
			if ok1 && ok2 {
				dst[k] = mergeMaps(m1, m2)
				continue
			}
			// keep existing scalar
			continue
		}
		dst[k] = v
	}
	return dst
}
