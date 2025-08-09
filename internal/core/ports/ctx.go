package ports

type Ctx interface {
	Values() map[string]any // .Name, .Module, etc.

	SetValue(key string, value any)

	ProjectRoot() string

	FS() FSWriter
	Renderer() Renderer

	GoMod() GoModEditor
}
