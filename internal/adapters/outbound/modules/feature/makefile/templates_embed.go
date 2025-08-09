package makefile

import "embed"

// TemplatesFS embeds all files under the local templates directory of the makefile module.
//
//go:embed templates
var TemplatesFS embed.FS
