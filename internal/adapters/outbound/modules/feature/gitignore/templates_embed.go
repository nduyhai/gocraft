package gitignore

import "embed"

// TemplatesFS embeds the template files for the gitignore module.
// Explicitly include .gitignore.tmpl since go:embed ignores dotfiles unless directly specified.
//
//go:embed templates/.gitignore.tmpl
var TemplatesFS embed.FS
