package entity

// Project represents the desired project to generate.
type Project struct {
	Name      string            // e.g., myapp
	Module    string            // e.g., github.com/you/myapp
	Template  string            // e.g., "basic"
	Options   map[string]string // arbitrary options for templates
	TargetDir string            // directory where to generate
}
