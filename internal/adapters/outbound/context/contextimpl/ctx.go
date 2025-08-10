package contextimpl

import "github.com/nduyhai/gocraft/internal/core/ports"

// Ctx implements ports.Ctx and holds outbound collaborators and generation values.
type Ctx struct {
	projectRoot    string
	values         map[string]any
	fs             ports.FSWriter
	renderer       ports.Renderer
	gomod          ports.GoModEditor
	adaptersModule ports.DependencyInjectionEditor
	config         ports.ConfigEditor
}

// New constructs a new Ctx.
// values may be nil; an internal map will be allocated.
func New(projectRoot string, fs ports.FSWriter, renderer ports.Renderer, gomod ports.GoModEditor, adapters ports.DependencyInjectionEditor, config ports.ConfigEditor, values map[string]any) *Ctx {
	if values == nil {
		values = make(map[string]any)
	}
	return &Ctx{
		projectRoot:    projectRoot,
		values:         values,
		fs:             fs,
		renderer:       renderer,
		gomod:          gomod,
		adaptersModule: adapters,
		config:         config,
	}
}

// Values returns the context values map used in templates and path tokens.
func (c *Ctx) Values() map[string]any { return c.values }

// SetValue sets a key/value into the context values map.
func (c *Ctx) SetValue(key string, value any) { c.values[key] = value }

// ProjectRoot returns the root directory for project generation.
func (c *Ctx) ProjectRoot() string { return c.projectRoot }

// FS returns the file system writer.
func (c *Ctx) FS() ports.FSWriter { return c.fs }

// Renderer returns the template renderer.
func (c *Ctx) Renderer() ports.Renderer { return c.renderer }

// GoMod returns the go.mod editor utility.
func (c *Ctx) GoMod() ports.GoModEditor { return c.gomod }

// AdaptersModule returns the adapters module file editor.
func (c *Ctx) AdaptersModule() ports.DependencyInjectionEditor { return c.adaptersModule }

// Config returns the config.yaml editor.
func (c *Ctx) Config() ports.ConfigEditor { return c.config }
