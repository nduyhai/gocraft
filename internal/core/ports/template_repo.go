package ports

// Template represents a named template repository entry that can be rendered.
type Template struct {
	Name  string // template name (e.g., "basic")
	Files []TmplFile
}

type TmplFile struct {
	Path    string
	Content string
}

type TemplateRepo interface {
	// Load returns the template by name.
	Load(name string) (Template, error)
	// Names lists available template names.
	Names() []string
}
